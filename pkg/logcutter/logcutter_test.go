package logcutter_test

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
	"fmt"
	"reflect"
	"strings"
	"testing"

	v1 "github.com/bhojpur/piro/pkg/api/v1"
	"github.com/bhojpur/piro/pkg/logcutter"
)

func TestDefaultCutterSlice(t *testing.T) {
	tests := []struct {
		Input  string
		Events []v1.LogSliceEvent
		Error  error
	}{
		{
			`
[foobar] Hello World this is a test
[otherproc] Some other process
[foobar] More output
[foobar|DONE]
[otherproc] Cool beans
			`,
			[]v1.LogSliceEvent{
				v1.LogSliceEvent{Name: "foobar", Type: v1.LogSliceType_SLICE_START},
				v1.LogSliceEvent{Name: "foobar", Type: v1.LogSliceType_SLICE_CONTENT, Payload: "Hello World this is a test"},
				v1.LogSliceEvent{Name: "otherproc", Type: v1.LogSliceType_SLICE_START},
				v1.LogSliceEvent{Name: "otherproc", Type: v1.LogSliceType_SLICE_CONTENT, Payload: "Some other process"},
				v1.LogSliceEvent{Name: "foobar", Type: v1.LogSliceType_SLICE_CONTENT, Payload: "More output"},
				v1.LogSliceEvent{Name: "foobar", Type: v1.LogSliceType_SLICE_DONE},
				v1.LogSliceEvent{Name: "otherproc", Type: v1.LogSliceType_SLICE_CONTENT, Payload: "Cool beans"},
				v1.LogSliceEvent{Name: "otherproc", Type: v1.LogSliceType_SLICE_ABANDONED},
			},
			nil,
		},
		{
			`
[build|PHASE] Pushing foobar
[components/foobar:docker] c13a632cd17b: Preparing
			`,
			[]v1.LogSliceEvent{
				v1.LogSliceEvent{Name: "build", Type: v1.LogSliceType_SLICE_PHASE, Payload: "Pushing foobar"},
				v1.LogSliceEvent{Name: "components/foobar:docker", Type: v1.LogSliceType_SLICE_START},
				v1.LogSliceEvent{Name: "components/foobar:docker", Type: v1.LogSliceType_SLICE_CONTENT, Payload: "c13a632cd17b: Preparing"},
				v1.LogSliceEvent{Name: "components/foobar:docker", Type: v1.LogSliceType_SLICE_ABANDONED},
			},
			nil,
		},
	}

	for _, test := range tests {
		content := strings.TrimSpace(test.Input)
		evtchan, errchan := logcutter.DefaultCutter.Slice(bytes.NewReader([]byte(content)))

		var (
			events []v1.LogSliceEvent
			err    error
		)
	recv:
		for {
			select {
			case evt := <-evtchan:
				if evt == nil {
					break recv
				}

				events = append(events, *evt)
			case err = <-errchan:
				break recv
			}
		}

		if err != test.Error {
			t.Errorf("unexpected error: \"%s\", expected \"%s\"", err, test.Error)
		}
		if !reflect.DeepEqual(test.Events, events) {
			expevt := make([]string, len(test.Events))
			for i, evt := range test.Events {
				expevt[i] = fmt.Sprintf("\t[%s] %s: %s", evt.Name, evt.Type.String(), evt.Payload)
			}
			actevt := make([]string, len(events))
			for i, evt := range events {
				actevt[i] = fmt.Sprintf("\t[%s] %s: %s", evt.Name, evt.Type.String(), evt.Payload)
			}

			t.Errorf("unexpected events:\n%s\nexpected:\n%s", strings.Join(actevt, "\n"), strings.Join(expevt, "\n"))
		}
	}
}
