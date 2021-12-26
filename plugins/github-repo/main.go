package main

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
	"github.com/google/go-github/v31/github"
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
