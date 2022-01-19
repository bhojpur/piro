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
	"reflect"

	v1 "github.com/bhojpur/piro/pkg/api/v1"
	plugin "github.com/bhojpur/piro/pkg/plugin/client"
	log "github.com/sirupsen/logrus"
)

// Config configures this plugin
type Config struct {
	Emoji string `yaml:"emoji"`
}

func main() {
	plugin.Serve(&Config{},
		plugin.WithIntegrationPlugin(&integrationPlugin{}),
	)
}

type integrationPlugin struct{}

func (*integrationPlugin) Run(ctx context.Context, config interface{}, srv v1.PiroServiceClient) error {
	cfg, ok := config.(*Config)
	if !ok {
		return fmt.Errorf("config has wrong type %s", reflect.TypeOf(config))
	}

	sub, err := srv.Subscribe(ctx, &v1.SubscribeRequest{})
	if err != nil {
		return err
	}

	log.Infof("hello world %s", cfg.Emoji)
	for {
		resp, err := sub.Recv()
		if err != nil {
			return err
		}

		fmt.Printf("%s %v\n", cfg.Emoji, resp)
	}
}
