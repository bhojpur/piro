package provider

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
	"bufio"
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/bhojpur/piro/pkg/plugin/common"

	"github.com/google/go-github/v42/github"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	corev1 "k8s.io/api/core/v1"
)

const (
	// defaultContainerImage is the image we use to clone the code, unless configured otherwise
	defaultContainerImage = "alpine/git:latest"
)

// GitCredentialHelper can authenticate provide authentication credentials for a repository
type GitCredentialHelper func(ctx context.Context) (user string, pass string, err error)

// GithubRepoServer provides access to Github repos
type GithubRepoServer struct {
	Client *github.Client
	Auth   GitCredentialHelper

	Config Config
	common.UnimplementedRepositoryPluginServer
}

// Config configures the GithubRepoServer
type Config struct {
	ContainerImage string
}

// RepoHost returns the host which this plugins integrates with
func (s *GithubRepoServer) RepoHost(context.Context, *common.RepoHostRequest) (*common.RepoHostResponse, error) {
	return &common.RepoHostResponse{
		Host: "github.com",
	}, nil
}

// Resolve resolves the repo's revision based on its ref(erence)
func (s *GithubRepoServer) Resolve(ctx context.Context, req *common.ResolveRequest) (*common.ResolveResponse, error) {
	repo := req.Repository
	if repo.Revision != "" {
		return &common.ResolveResponse{Repository: repo}, nil
	}
	if repo.Ref == "" {
		return nil, status.Errorf(codes.InvalidArgument, "ref is empty")
	}

	branch, _, err := s.Client.Repositories.GetBranch(ctx, repo.Owner, repo.Repo, repo.Ref, false)
	if err != nil {
		return nil, translateGitHubToGRPCError(err, repo.Revision, repo.Ref)
	}
	if branch == nil {
		return nil, status.Error(codes.NotFound, "did not find ref")
	}
	if branch.Commit == nil || branch.Commit.SHA == nil {
		return nil, status.Error(codes.NotFound, "ref did not point to a commit")
	}
	repo.Revision = *branch.Commit.SHA
	log.WithField("ref", repo.Ref).WithField("rev", repo.Revision).Debug("resolved reference to revision")

	return &common.ResolveResponse{
		Repository: repo,
	}, nil
}

// GetRemoteAnnotations extracts Bhojpur Piro annotations form information associated
// with a particular commit, e.g. the commit message, Pull Requests or merge requests.
func (s *GithubRepoServer) GetRemoteAnnotations(ctx context.Context, req *common.GetRemoteAnnotationsRequest) (resp *common.GetRemoteAnnotationsResponse, err error) {
	repo := req.Repository

	opts := &github.ListOptions{Page: 2, PerPage: 2}
	commit, _, err := s.Client.Repositories.GetCommit(ctx, repo.Owner, repo.Repo, repo.Revision, opts)
	if err != nil {
		return nil, translateGitHubToGRPCError(err, repo.Revision, repo.Ref)
	}

	res := make(map[string]string)
	if commit.Commit != nil {
		atns := ParseAnnotations(commit.Commit.GetMessage())
		for k, v := range atns {
			res[k] = v
		}
	}

	prs, _, err := s.Client.PullRequests.ListPullRequestsWithCommit(ctx, repo.Owner, repo.Repo, commit.GetSHA(), &github.PullRequestListOptions{
		State: "open",
		Sort:  "created",
	})
	if err != nil {
		log.WithField("ref", repo.Ref).WithField("rev", repo.Revision).WithError(err).Warn("cannot search for associated PR's")
	} else if len(prs) >= 1 {
		sort.Slice(prs, func(i, j int) bool { return prs[i].GetCreatedAt().Before(prs[j].GetCreatedAt()) })
		pr := prs[0]

		if len(prs) > 1 {
			log.WithField("ref", repo.Ref).WithField("rev", repo.Revision).WithField("pr", pr.GetHTMLURL()).Warn("found multiple open PR's - using oldest one")
		}

		atns := ParseAnnotations(pr.GetBody())
		for k, v := range atns {
			res[k] = v
		}
	}
	return &common.GetRemoteAnnotationsResponse{
		Annotations: res,
	}, nil
}

// PareseAnnotations parses one annotation per line in the form of "/piro <key>(=<value>)?".
// Any line not matching this format is silently ignored.
func ParseAnnotations(message string) (res map[string]string) {
	scanner := bufio.NewScanner(bytes.NewReader([]byte(message)))
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, "- [x]")
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "/piro ") {
			continue
		}
		line = strings.TrimSpace(strings.TrimPrefix(line, "/piro "))

		var name, val string
		if idx := strings.Index(line, "="); idx > -1 {
			name = line[:idx]
			val = strings.TrimPrefix(line[idx:], "=")
		} else {
			name = line
			val = ""
		}
		if res == nil {
			res = make(map[string]string)
		}
		res[name] = val
	}
	return res
}

// ContentInitContainer produces the init container YAML required to initialize
// the build context from this repository in /application.
func (s *GithubRepoServer) ContentInitContainer(ctx context.Context, req *common.ContentInitContainerRequest) (*common.ContentInitContainerResponse, error) {
	image := s.Config.ContainerImage
	if image == "" {
		image = defaultContainerImage
	}

	var (
		repo = req.Repository
		user string
		pass string
	)
	if s.Auth != nil {
		var err error
		user, pass, err = s.Auth(context.Background())
		if err != nil {
			return nil, err
		}
	}

	cloneCmd := "git clone"
	if user != "" || pass != "" {
		cloneCmd = fmt.Sprintf("git clone -c \"credential.helper=/bin/sh -c 'echo username=$GHUSER_SECRET; echo password=$GHPASS_SECRET'\"")
	}
	cloneCmd = fmt.Sprintf("%s https://github.com/%s/%s.git .; git checkout %s", cloneCmd, repo.Owner, repo.Repo, repo.Revision)

	c := []corev1.Container{
		{
			Name:  "github-checkout",
			Image: image,
			Command: []string{
				"sh", "-c",
				cloneCmd,
			},
			Env: []corev1.EnvVar{
				{
					Name:  "GHUSER_SECRET",
					Value: user,
				},
				{
					Name:  "GHPASS_SECRET",
					Value: pass,
				},
			},
			WorkingDir: "/application",
		},
	}
	buf := bytes.NewBuffer(nil)
	err := gob.NewEncoder(buf).Encode(c)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot serialize init container: %q", err)
	}

	return &common.ContentInitContainerResponse{
		Container: buf.Bytes(),
	}, nil
}

// Download downloads a file from the repository.
func (s *GithubRepoServer) Download(ctx context.Context, req *common.DownloadRequest) (*common.DownloadResponse, error) {
	dl, resp, err := s.Client.Repositories.DownloadContents(ctx, req.Repository.Owner, req.Repository.Repo, req.Path, &github.RepositoryContentGetOptions{
		Ref: req.Repository.Revision,
	})
	if err != nil {
		return nil, translateGitHubToGRPCError(err, req.Repository.Revision, req.Repository.Ref)
	}
	defer dl.Close()

	res, err := ioutil.ReadAll(dl)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot download file: %q", err)
	}
	if resp == nil {
		return nil, status.Errorf(codes.Internal, "cannot download nil resp file: %q", nil)
	}
	return &common.DownloadResponse{
		Content: res,
	}, nil
}

// ListFiles lists all files in a directory.
func (s *GithubRepoServer) ListFiles(ctx context.Context, req *common.ListFilesRequest) (*common.ListFilesReponse, error) {
	_, dc, _, err := s.Client.Repositories.GetContents(ctx, req.Repository.Owner, req.Repository.Repo, req.Path, &github.RepositoryContentGetOptions{
		Ref: req.Repository.Revision,
	})
	if err != nil {
		return nil, translateGitHubToGRPCError(err, req.Repository.Revision, req.Repository.Ref)
	}

	res := make([]string, 0, len(dc))
	for _, cnt := range dc {
		if cnt.GetType() != "file" {
			continue
		}
		res = append(res, cnt.GetPath())
	}
	return &common.ListFilesReponse{Paths: res}, nil
}

func translateGitHubToGRPCError(err error, rev, ref string) error {
	if gherr, ok := err.(*github.ErrorResponse); ok && gherr.Response.StatusCode == 422 {
		msg := fmt.Sprintf("revision %s", rev)
		if ref != "" {
			msg = fmt.Sprintf("ref %s (revision %s)", ref, rev)
		}
		return status.Error(codes.NotFound, fmt.Sprintf("%s not found", msg))
	}

	return status.Error(codes.Internal, err.Error())
}
