package store

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
	"io"
	"io/ioutil"
	"sync"

	v1 "github.com/bhojpur/piro/pkg/api/v1"
	"github.com/bhojpur/piro/pkg/filterexpr"
	"golang.org/x/xerrors"
)

// NewInMemoryLogStore provides a new log store which stores its logs in memory
func NewInMemoryLogStore() Logs {
	return &inMemoryLogStore{
		logs: make(map[string]*logSession),
	}
}

// inMemoryLogStore implements a log store in memory
type inMemoryLogStore struct {
	logs map[string]*logSession
	mu   sync.RWMutex
}

type logSession struct {
	Data   *bytes.Buffer
	Reader map[chan []byte]struct{}
	Mu     sync.RWMutex
}

func (l *logSession) Write(p []byte) (n int, err error) {
	l.Mu.Lock()
	defer l.Mu.Unlock()

	n, err = l.Data.Write(p)
	if n > 0 {

		for r := range l.Reader {
			r <- p[:n]
		}
	}
	if err != nil {
		return n, err
	}
	return
}

func (l *logSession) Close() error {
	return nil
}

type logSessionReader struct {
	Log       *logSession
	Pos       int
	R         chan []byte
	remainder []byte
	closed    bool
}

func (lr *logSessionReader) Read(p []byte) (n int, err error) {
	if lr.closed {
		return 0, io.ErrClosedPipe
	}

	if len(lr.remainder) > 0 {
		n = copy(p, lr.remainder)
		lr.remainder = lr.remainder[:n]
		lr.Pos += n
		return
	}

	lr.Log.Mu.RLock()
	if lr.Pos >= lr.Log.Data.Len() {
		lr.Log.Mu.RUnlock()
		inc := <-lr.R

		n = copy(p, inc)
		lr.remainder = inc[:n]
		lr.Pos += n
		return
	}

	n = copy(p, lr.Log.Data.Bytes()[lr.Pos:])
	lr.Pos += n

	lr.Log.Mu.RUnlock()
	return 0, nil
}

func (lr *logSessionReader) Close() error {
	lr.Log.Mu.Lock()
	defer lr.Log.Mu.Unlock()

	delete(lr.Log.Reader, lr.R)
	lr.closed = true

	return nil
}

// Place writes to this store
func (s *inMemoryLogStore) Open(id string) (io.WriteCloser, error) {
	s.mu.Lock()
	if _, ok := s.logs[id]; ok {
		s.mu.Unlock()
		return nil, ErrAlreadyExists
	}

	lg := &logSession{
		Data:   bytes.NewBuffer(nil),
		Reader: make(map[chan []byte]struct{}),
	}

	s.logs[id] = lg
	s.mu.Unlock()

	return lg, nil
}

func (s *inMemoryLogStore) Write(id string) (io.Writer, error) {
	return nil, xerrors.Errorf("not supported")
}

// Read reads from this store
func (s *inMemoryLogStore) Read(id string) (io.ReadCloser, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	l, ok := s.logs[id]
	if !ok {
		return nil, ErrNotFound
	}

	ch := make(chan []byte)
	l.Mu.Lock()
	l.Reader[ch] = struct{}{}
	l.Mu.Unlock()
	return ioutil.NopCloser(&logSessionReader{
		Log: l,
		R:   ch,
	}), nil
}

// NewInMemoryJobStore creates a new in-memory Job store
func NewInMemoryJobStore() Jobs {
	return &inMemoryJobStore{
		jobs:  make(map[string]v1.JobStatus),
		specs: make(map[string][]byte),
	}
}

type inMemoryJobStore struct {
	jobs  map[string]v1.JobStatus
	specs map[string][]byte
	mu    sync.RWMutex
}

// Store stores Job information in the store.
// Storing a job whose name we already have in store will override the previously
// stored job.
func (s *inMemoryJobStore) Store(ctx context.Context, job v1.JobStatus) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.jobs[job.Name] = job
	return nil
}

// Retrieves a particular job bassd on its name.
// If the Job is unknown we'll return ErrNotFound.
func (s *inMemoryJobStore) Get(ctx context.Context, name string) (*v1.JobStatus, error) {
	s.mu.RLock()
	job, ok := s.jobs[name]
	s.mu.RUnlock()

	if !ok {
		return nil, ErrNotFound
	}

	return &job, nil
}

// Searches for Job(s) based on their annotations
func (s *inMemoryJobStore) Find(ctx context.Context, filter []*v1.FilterExpression, order []*v1.OrderExpression, start, limit int) (slice []v1.JobStatus, total int, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var res []v1.JobStatus
	for _, js := range s.jobs {
		if !filterexpr.MatchesFilter(&js, filter) {
			continue
		}
		res = append(res, js)
	}
	return res, len(res), nil
}

func (s *inMemoryJobStore) StoreJobSpec(name string, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.specs[name] = data
	return nil
}

func (s *inMemoryJobStore) GetJobSpec(name string) (data []byte, err error) {
	s.mu.RLock()
	s.mu.RUnlock()

	data, ok := s.specs[name]
	if !ok {
		return nil, ErrNotFound
	}
	return data, nil
}
