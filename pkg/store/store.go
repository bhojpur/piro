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
	"context"
	"fmt"
	"io"

	v1 "github.com/bhojpur/piro/pkg/api/v1"
)

var (
	// ErrNotFound is returned by Read if something isn't found
	ErrNotFound = fmt.Errorf("not found")

	// ErrAlreadyExists is returned when attempting to place something which already exists
	ErrAlreadyExists = fmt.Errorf("exists already")
)

// Logs provides access to the logstore
type Logs interface {
	// Open places a logfile in this store.
	// The caller is expected to close this writer when the task is complete.
	// If the logfile is already open we'll return an error.
	Open(id string) (io.WriteCloser, error)

	// Write writes to a previously placed logfile.
	// If the logfile is unknown, we'll return an error.
	Write(id string) (io.Writer, error)

	// Read retrieves a log file from this store.
	// Returns ErrNotFound if the log file isn't found.
	// Callers are supposed to close the reader once done.
	// Reading from logs currently being written is supported.
	Read(id string) (io.ReadCloser, error)
}

// Jobs provides access to past Jobs
type Jobs interface {
	// Store stores schedulable Kubernetes Job information in the Memory / DB store.
	// Storing a Job whose name we already have in store will override the previously
	// stored Job.
	Store(ctx context.Context, job v1.JobStatus) error

	// StoreJobSpec stores Job YAML data.
	StoreJobSpec(name string, data []byte) error

	// Retrieves a particular Kubernetes Job bassd on its name.
	// If the Kubernetes Job is unknown we'll return ErrNotFound.
	Get(ctx context.Context, name string) (*v1.JobStatus, error)

	// Get retrieves previously stored Kubernetes Job spec data
	GetJobSpec(name string) (data []byte, err error)

	// Searches for Kubernetes Jobs based on their annotations. If filter is
	// empty no filter is applied. If limit is 0, no limit is applied.
	Find(ctx context.Context, filter []*v1.FilterExpression, order []*v1.OrderExpression, start, limit int) (slice []v1.JobStatus, total int, err error)
}

// NumberGroup enables to atomic generation and storage of numbers.
// This is used for build numbering
type NumberGroup interface {
	// Latest returns the latest number of a particular number group.
	// Returns ErrNotFound if the group does not exist. A zero result is a valid
	// number in a group and does not indicate its non-existence.
	Latest(group string) (nr int, err error)

	// Next returns the next number in the group. If the group did not exist prior
	// to this call it is created. This function is thread-safe and atomic.
	Next(group string) (nr int, err error)
}
