package host

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
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"io"
	"io/ioutil"
	"sync"
	"time"

	v1 "github.com/bhojpur/piro/pkg/api/v1"
	piro "github.com/bhojpur/piro/pkg/engine"
	"github.com/bhojpur/piro/pkg/plugin/common"

	corev1 "k8s.io/api/core/v1"
)

type compoundRepositoryProvider struct {
	hosts map[string]piro.RepositoryProvider
	mu    sync.RWMutex
}

var errNoRepoProvider = errors.New("no host provider")

func (c *compoundRepositoryProvider) getProvider(host string) (piro.RepositoryProvider, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	res, ok := c.hosts[host]
	if !ok {
		return nil, errNoRepoProvider
	}
	return res, nil
}

func (c *compoundRepositoryProvider) registerProvider(host string, prov piro.RepositoryProvider) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.hosts == nil {
		c.hosts = make(map[string]piro.RepositoryProvider)
	}

	c.hosts[host] = prov
}

// Resolve resolves the repo's revision based on its ref(erence).
// If the revision is already set, this operation does nothing.
func (c *compoundRepositoryProvider) Resolve(ctx context.Context, repo *v1.Repository) error {
	prov, err := c.getProvider(repo.Host)
	if err != nil {
		return err
	}
	return prov.Resolve(ctx, repo)
}

// RemoteAnnotations extracts Bhojpur Piro annotations form information associated
// with a particular commit, e.g. the commit message, PRs or merge requests.
// Implementors can expect the revision of the repo object to be set.
func (c *compoundRepositoryProvider) RemoteAnnotations(ctx context.Context, repo *v1.Repository) (annotations map[string]string, err error) {
	prov, err := c.getProvider(repo.Host)
	if err != nil {
		return nil, err
	}
	return prov.RemoteAnnotations(ctx, repo)
}

// ContentProvider produces a content provider for a particular repo
func (c *compoundRepositoryProvider) ContentProvider(ctx context.Context, repo *v1.Repository) (piro.ContentProvider, error) {
	prov, err := c.getProvider(repo.Host)
	if err != nil {
		return nil, err
	}
	return prov.ContentProvider(ctx, repo)
}

// FileProvider provides direct access to repository content
func (c *compoundRepositoryProvider) FileProvider(ctx context.Context, repo *v1.Repository) (piro.FileProvider, error) {
	prov, err := c.getProvider(repo.Host)
	if err != nil {
		return nil, err
	}
	return prov.FileProvider(ctx, repo)
}

type pluginHostProvider struct {
	C common.RepositoryPluginClient
}

var _ piro.RepositoryProvider = &pluginHostProvider{}
var _ piro.FileProvider = &pluginContentProvider{}

// Resolve resolves the repo's revision based on its ref(erence).
// If the revision is already set, this operation does nothing.
func (p *pluginHostProvider) Resolve(ctx context.Context, repo *v1.Repository) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := p.C.Resolve(ctx, &common.ResolveRequest{
		Repository: repo,
	})
	if err != nil {
		return err
	}
	*repo = *resp.Repository
	return nil
}

// RemoteAnnotations extracts Bhojpur Piro annotations form information associated
// with a particular commit, e.g. the commit message, PRs or merge requests.
func (p *pluginHostProvider) RemoteAnnotations(ctx context.Context, repo *v1.Repository) (annotations map[string]string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := p.C.GetRemoteAnnotations(ctx, &common.GetRemoteAnnotationsRequest{Repository: repo})
	if err != nil {
		return nil, err
	}

	return resp.Annotations, nil
}

// ContentProvider produces a content provider for a particular repo
func (p *pluginHostProvider) ContentProvider(ctx context.Context, repo *v1.Repository) (piro.ContentProvider, error) {
	return &pluginContentProvider{
		Repo: repo,
		C:    p.C,
	}, nil
}

// FileProvider provides direct access to repository content
func (p *pluginHostProvider) FileProvider(ctx context.Context, repo *v1.Repository) (piro.FileProvider, error) {
	return &pluginContentProvider{
		Repo: repo,
		C:    p.C,
	}, nil
}

type pluginContentProvider struct {
	Repo *v1.Repository
	C    common.RepositoryPluginClient
}

func (c *pluginContentProvider) InitContainer() (res []corev1.Container, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.C.ContentInitContainer(ctx, &common.ContentInitContainerRequest{
		Repository: c.Repo,
	})
	if err != nil {
		return nil, err
	}

	err = gob.NewDecoder(bytes.NewReader(resp.Container)).Decode(&res)
	return
}

func (c *pluginContentProvider) Serve(jobName string) error {
	return nil
}

func (c *pluginContentProvider) Download(ctx context.Context, path string) (io.ReadCloser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.C.Download(ctx, &common.DownloadRequest{
		Repository: c.Repo,
		Path:       path,
	})
	if err != nil {
		return nil, err
	}

	return ioutil.NopCloser(bytes.NewReader(resp.Content)), nil
}

func (c *pluginContentProvider) ListFiles(ctx context.Context, path string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.C.ListFiles(ctx, &common.ListFilesRequest{
		Repository: c.Repo,
		Path:       path,
	})
	if err != nil {
		return nil, err
	}

	return resp.Paths, nil
}
