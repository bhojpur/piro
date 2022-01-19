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
	"strings"

	v1 "github.com/bhojpur/piro/pkg/api/v1"
	plugin "github.com/bhojpur/piro/pkg/plugin/client"
	"github.com/bhojpur/piro/pkg/reporef"
	cron "github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

// Config configures this plugin
type Config struct {
	Tasks []struct {
		Spec        string            `yaml:"spec"`
		Repo        string            `yaml:"repo"`
		JobPath     string            `yaml:"jobPath,omitempty"`
		Trigger     string            `yaml:"trigger,omitempty"`
		Annotations map[string]string `yaml:"annotations,omitempty"`
	} `yaml:"tasks"`
}

func main() {
	plugin.Serve(&Config{},
		plugin.WithIntegrationPlugin(&cronPlugin{}),
	)
}

type cronPlugin struct{}

func (*cronPlugin) Run(ctx context.Context, config interface{}, srv v1.PiroServiceClient) error {
	cfg, ok := config.(*Config)
	if !ok {
		return fmt.Errorf("config has wrong type %s", reflect.TypeOf(config))
	}

	c := cron.New()
	for idx, task := range cfg.Tasks {
		repo, err := reporef.Parse(task.Repo)
		if err != nil {
			return err
		}

		var trigger v1.JobTrigger
		if trg, ok := v1.JobTrigger_value[fmt.Sprintf("TRIGGER_%s", strings.ToUpper(task.Trigger))]; ok {
			trigger = v1.JobTrigger(trg)
		} else if task.Trigger != "" {
			return fmt.Errorf("unknown Job trigger %s", task.Trigger)
		}

		var annotations []*v1.Annotation
		for k, v := range task.Annotations {
			annotations = append(annotations, &v1.Annotation{
				Key:   k,
				Value: v,
			})
		}

		request := &v1.StartGitHubJobRequest{
			Metadata: &v1.JobMetadata{
				Owner:       "cron",
				Annotations: annotations,
				Trigger:     trigger,
				Repository:  repo,
			},
			JobPath: task.JobPath,
		}
		_ = c.AddFunc(task.Spec, func() {
			_, err := srv.StartGitHubJob(ctx, request)
			if err != nil {
				log.WithError(err).WithField("idx", idx).WithField("spec", task.Spec).Error("cannot start Job")
			}
		})
		if err != nil {
			return err
		}

		log.WithField("spec", task.Spec).Info("scheduled job")
	}
	c.Run()

	return nil
}
