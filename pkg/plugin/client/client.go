package client

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"context"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	v1 "github.com/bhojpur/piro/pkg/api/v1"
	"github.com/bhojpur/piro/pkg/plugin/common"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
)

type Services struct {
	v1.PiroServiceClient
	v1.PiroUIClient
}

// IntegrationPlugin works on the public Bhojpur Piro API
type IntegrationPlugin interface {
	// Run runs the plugin. Once this function returns the plugin stops running.
	// Implementors must respect the context deadline as that's the signal for graceful shutdown.
	Run(ctx context.Context, config interface{}, srv *Services) error
}

// RepositoryPlugin adds support for a repository host
type RepositoryPlugin interface {
	// Run runs the plugin. The plugin runs until the context is canceled and the server returned
	// by this function is expected to remain functional until then.
	Run(ctx context.Context, config interface{}) (common.RepositoryPluginServer, error)
}

// AuthenticationPlugin adds support for API authentication
type AuthenticationPlugin interface {
	// Run runs the plugin. The plugin runs until the context is canceled and the server returned
	// by this function is expected to remain functional until then.
	Run(ctx context.Context, config interface{}) (common.AuthenticationPluginServer, error)
}

// ServeOpt configures a plugin serve
type ServeOpt struct {
	Type common.Type
	Run  func(ctx context.Context, config interface{}, socket string) error
}

// WithIntegrationPlugin registers integration plugin capabilities
func WithIntegrationPlugin(p IntegrationPlugin) ServeOpt {
	return ServeOpt{
		Type: common.TypeIntegration,
		Run: func(ctx context.Context, config interface{}, socket string) error {
			conn, err := grpc.Dial(socket, grpc.WithInsecure(), grpc.WithDialer(unixConnect))
			if err != nil {
				return xerrors.Errorf("did not connect: %v", err)
			}
			defer conn.Close()

			return p.Run(ctx, config, &Services{
				PiroServiceClient: v1.NewPiroServiceClient(conn),
				PiroUIClient:      v1.NewPiroUIClient(conn),
			})
		},
	}
}

// WithRepositoryPlugin registers repo plugin capabilities
func WithRepositoryPlugin(p RepositoryPlugin) ServeOpt {
	return ServeOpt{
		Type: common.TypeRepository,
		Run: func(ctx context.Context, config interface{}, socket string) error {
			lis, err := net.Listen("unix", socket)
			if err != nil {
				return err
			}
			service, err := p.Run(ctx, config)
			if err != nil {
				return err
			}

			s := grpc.NewServer()
			common.RegisterRepositoryPluginServer(s, service)
			return s.Serve(lis)
		},
	}
}

// WithRepositoryPlugin registers repo plugin capabilities
func WithAuthenticationPlugin(p AuthenticationPlugin) ServeOpt {
	return ServeOpt{
		Type: common.TypeAuthentication,
		Run: func(ctx context.Context, config interface{}, socket string) error {
			lis, err := net.Listen("unix", socket)
			if err != nil {
				return err
			}
			service, err := p.Run(ctx, config)
			if err != nil {
				return err
			}

			s := grpc.NewServer()
			common.RegisterAuthenticationPluginServer(s, service)
			return s.Serve(lis)
		},
	}
}

const proxyPassPluginType common.Type = "proxy-pass"

// ProxyPassPlugin adds additional support for proxied webhooks
type ProxyPassPlugin interface {
	Serve(ctx context.Context, l net.Listener) error
}

// WithProxyPass enables a "through Piro" proxy route to the plugin.
// The route will be available at "http://<piro-location>/plugins/<plugin-name>"
func WithProxyPass(p ProxyPassPlugin) ServeOpt {
	return ServeOpt{
		Type: proxyPassPluginType,
		Run: func(ctx context.Context, config interface{}, socket string) error {
			l, err := net.Listen("unix", socket)
			if err != nil {
				return err
			}
			return p.Serve(ctx, l)
		},
	}
}

// Serve is the main entry point for plugins
func Serve(configType interface{}, opts ...ServeOpt) {
	if typ := reflect.TypeOf(configType); typ.Kind() != reflect.Ptr {
		log.Fatal("configType is not a pointer")
	}

	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	errchan := make(chan error)

	if len(os.Args) != 4 {
		log.Fatalf("usage: %s <type> <cfgfile.yaml> <socket>", os.Args[0])
		return
	}
	tpe, cfgfn, socketfn := os.Args[1], os.Args[2], os.Args[3]

	// load config
	cfgraw, err := ioutil.ReadFile(cfgfn)
	if err != nil {
		log.Fatalf("cannot read config file: %v", err)
	}

	err = yaml.Unmarshal(cfgraw, configType)
	if err != nil {
		log.Fatalf("cannot unmarshal config: %v", err)
	}
	config := configType

	ctx, cancel := context.WithCancel(context.Background())
	for _, o := range opts {
		if o.Type != proxyPassPluginType {
			continue
		}

		o := o
		go func() {
			err = o.Run(ctx, nil, os.Getenv("PIRO_PLUGIN_PROXY_SOCKET"))
			if err != nil && err != context.Canceled {
				errchan <- err
			}
		}()
		break
	}

	var sv *ServeOpt
	for _, o := range opts {
		if string(o.Type) == tpe {
			sv = &o
			break
		}
	}
	if sv == nil {
		log.Fatalf("cannot serve as %s plugin", tpe)
	}
	go func() {
		err := sv.Run(ctx, config, socketfn)
		if err != nil && err != context.Canceled {
			errchan <- err
		}
	}()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)
	log.Info("plugin is running")
	select {
	case <-sigchan:
	case err := <-errchan:
		log.Fatal(err)
	}

	cancel()
}

func unixConnect(addr string, t time.Duration) (net.Conn, error) {
	return net.Dial("unix", addr)
}
