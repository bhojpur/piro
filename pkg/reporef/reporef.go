package reporef

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
	"strings"

	v1 "github.com/bhojpur/piro/pkg/api/v1"
	"golang.org/x/xerrors"
)

// Parse interprets a string pointing to a (GitHub) repository.
// We expect the string to be in the form of:
//    (host)/owner/repo(:ref|@sha)
func Parse(spec string) (*v1.Repository, error) {
	if strings.Contains(spec, ":") {
		segs := strings.Split(spec, ":")
		rep, ref := segs[0], segs[1]
		repo, err := parseRep(rep)
		if err != nil {
			return nil, err
		}
		repo.Ref = ref
		return repo, nil
	}
	if strings.Contains(spec, "@") {
		segs := strings.Split(spec, "@")
		rep, rev := segs[0], segs[1]
		repo, err := parseRep(rep)
		if err != nil {
			return nil, err
		}
		repo.Revision = rev
		return repo, nil
	}
	return parseRep(spec)
}

func parseRep(rep string) (*v1.Repository, error) {
	segs := strings.Split(rep, "/")
	if len(segs) < 2 || len(segs) > 3 {
		return nil, xerrors.Errorf("invalid repository spec")
	}

	res := &v1.Repository{}
	if len(segs) == 3 {
		res.Host = segs[0]
		segs = segs[1:]
	}
	res.Owner = segs[0]
	res.Repo = segs[1]
	return res, nil
}
