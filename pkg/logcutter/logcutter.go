package logcutter

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
	"io"
	"strings"

	v1 "github.com/bhojpur/piro/pkg/api/v1"
)

// Cutter splits a log stream into slices for more structured display
type Cutter interface {
	// Slice reads on the in reader line-by-line. For each line it can produce several events
	// on the events channel. Once the reader returns EOF the events and errchan are closed.
	// If anything goes wrong while reading a single error is written to errchan, but nothing is closed.
	Slice(in io.Reader) (events <-chan *v1.LogSliceEvent, errchan <-chan error)
}

const (
	// DefaultSlice is the parent slice of all unmarked content
	DefaultSlice = "default"
)

// NoCutter does not slice the content up at all
var NoCutter Cutter = noCutter{}

type noCutter struct{}

// Slice returns all log lines
func (noCutter) Slice(in io.Reader) (events <-chan *v1.LogSliceEvent, errchan <-chan error) {
	evts := make(chan *v1.LogSliceEvent)
	errc := make(chan error)
	events, errchan = evts, errc

	scanner := bufio.NewScanner(in)
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			evts <- &v1.LogSliceEvent{
				Name:    DefaultSlice,
				Type:    v1.LogSliceType_SLICE_CONTENT,
				Payload: line + "\n",
			}
		}
		if err := scanner.Err(); err != nil {
			errc <- err
		}
		close(evts)
		close(errc)
	}()

	return
}

// DefaultCutter implements the default cutting behaviour
var DefaultCutter Cutter = defaultCutter{}

type defaultCutter struct{}

// Slice cuts a log stream into pieces based on a configurable delimiter
func (defaultCutter) Slice(in io.Reader) (events <-chan *v1.LogSliceEvent, errchan <-chan error) {
	evts := make(chan *v1.LogSliceEvent)
	errc := make(chan error)
	events, errchan = evts, errc

	scanner := bufio.NewScanner(in)
	phase := DefaultSlice
	go func() {
		idx := make(map[string]struct{})
		for scanner.Scan() {
			line := scanner.Text()
			sl := strings.TrimSpace(line)

			var (
				name    string
				verb    string
				payload string
			)

			if !(strings.HasPrefix(sl, "[") && strings.Contains(sl, "]")) {
				name = phase
				payload = line
			} else {
				start := strings.IndexRune(sl, '[')
				end := strings.IndexRune(sl, ']')
				name = sl[start+1 : end]
				payload = strings.TrimPrefix(sl[end+1:], " ")

				if segs := strings.Split(name, "|"); len(segs) == 2 {
					name = segs[0]
					verb = segs[1]
				}
			}

			switch verb {
			case "DONE":
				delete(idx, name)
				evts <- &v1.LogSliceEvent{
					Name: name,
					Type: v1.LogSliceType_SLICE_DONE,
				}
				continue
			case "FAIL":
				delete(idx, name)
				evts <- &v1.LogSliceEvent{
					Name:    name,
					Payload: payload,
					Type:    v1.LogSliceType_SLICE_FAIL,
				}
				continue
			case "RESULT":
				evts <- &v1.LogSliceEvent{
					Name:    name,
					Type:    v1.LogSliceType_SLICE_RESULT,
					Payload: payload,
				}
				continue
			case "PHASE":
				evts <- &v1.LogSliceEvent{
					Name:    name,
					Type:    v1.LogSliceType_SLICE_PHASE,
					Payload: payload,
				}
				phase = name
				continue
			}

			_, exists := idx[name]
			if !exists {
				idx[name] = struct{}{}
				evts <- &v1.LogSliceEvent{
					Name: name,
					Type: v1.LogSliceType_SLICE_START,
				}
			}
			evts <- &v1.LogSliceEvent{
				Name:    name,
				Type:    v1.LogSliceType_SLICE_CONTENT,
				Payload: string([]byte(payload)),
			}
		}
		if err := scanner.Err(); err != nil {
			errc <- err
		}

		for name := range idx {
			evts <- &v1.LogSliceEvent{
				Name: name,
				Type: v1.LogSliceType_SLICE_ABANDONED,
			}
		}

		close(evts)
		close(errc)
	}()

	return
}
