package main

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
	"fmt"
	"net/http"
	"os"
	"reflect"

	plugin "github.com/bhojpur/piro/pkg/plugin/client"
	"github.com/bhojpur/piro/pkg/plugin/common"
	"github.com/bhojpur/piro/plugins/github-repo/pkg/provider"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v42/github"
)

// Config configures this plugin
type Config struct {
	PrivateKeyPath string `yaml:"privateKeyPath"`
	InstallationID int64  `yaml:"installationID,omitempty"`
	AppID          int64  `yaml:"appID"`

	ContainerImage string `yaml:"containerImage"`
}

func main() {
	plugin.Serve(&Config{},
		plugin.WithRepositoryPlugin(&githubRepoPlugin{}),
	)
	fmt.Fprintln(os.Stderr, "shutting down")
}

type githubRepoPlugin struct{}

func (*githubRepoPlugin) Run(ctx context.Context, config interface{}) (common.RepositoryPluginServer, error) {
	cfg, ok := config.(*Config)
	if !ok {
		return nil, fmt.Errorf("config has wrong type %s", reflect.TypeOf(config))
	}

	ghtr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, cfg.AppID, cfg.InstallationID, cfg.PrivateKeyPath)
	if err != nil {
		return nil, err
	}
	ghClient := github.NewClient(&http.Client{Transport: ghtr})

	return &provider.GithubRepoServer{
		Client: ghClient,
		Auth: func(ctx context.Context) (user string, pass string, err error) {
			tkn, err := ghtr.Token(ctx)
			if err != nil {
				return
			}
			user = "x-access-token"
			pass = tkn
			return
		},
		Config: provider.Config{
			ContainerImage: cfg.ContainerImage,
		},
	}, nil
}
