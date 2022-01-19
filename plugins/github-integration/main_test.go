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
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseCommand(t *testing.T) {
	type Expectation struct {
		Cmd  string
		Args []string
		Err  string
	}
	tests := []struct {
		Name        string
		Input       string
		Expectation Expectation
	}{
		{Name: "empty line", Input: ""},
		{Name: "ignore line", Input: "something\nsomethingelse"},
		{Name: "no command", Input: "/piro", Expectation: Expectation{Err: "missing command"}},
		{Name: "no arg", Input: "/piro foo", Expectation: Expectation{Cmd: "foo", Args: []string{}}},
		{Name: "one arg", Input: "/piro foo bar", Expectation: Expectation{Cmd: "foo", Args: []string{"bar"}}},
		{Name: "two args", Input: "/piro foo bar=baz something", Expectation: Expectation{Cmd: "foo", Args: []string{"bar=baz", "something"}}},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var (
				act Expectation
				err error
			)
			act.Cmd, act.Args, err = parseCommand(test.Input)
			if err != nil {
				act.Err = err.Error()
			}

			if diff := cmp.Diff(test.Expectation, act); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
