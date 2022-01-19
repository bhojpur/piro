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
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseAnnotations(t *testing.T) {
	tests := []struct {
		Name     string
		Input    string
		Expected map[string]string
	}{
		{
			Name:  "empty string",
			Input: "",
		},
		{
			Name:  "unrelated content",
			Input: "Something unrelated",
		},
		{
			Name:     "piro annotation",
			Input:    "/piro foobar",
			Expected: map[string]string{"foobar": ""},
		},
		{
			Name:     "piro annotation with value",
			Input:    "/piro foobar=value",
			Expected: map[string]string{"foobar": "value"},
		},
		{
			Name:     "piro annotation with checkbox",
			Input:    "- [x] /piro foobar",
			Expected: map[string]string{"foobar": ""},
		},
		{
			Name:     "piro annotation with checkbox",
			Input:    "- [x]    /piro foobar=value",
			Expected: map[string]string{"foobar": "value"},
		},
		{
			Name:  "piro annotation with unchecked list checkbox",
			Input: "- [ ] /piro foobar",
		},
		{
			Name:     "mixed piro annotation",
			Input:    "hello world\n  /piro foo=bar",
			Expected: map[string]string{"foo": "bar"},
		},
		{
			Name:     "piro annotation with complex value",
			Input:    "/piro foobar=this=is=another/value 12,3,4,5",
			Expected: map[string]string{"foobar": "this=is=another/value 12,3,4,5"},
		},
		{
			Name:     "piro annotation with empty value",
			Input:    "/piro foobar=",
			Expected: map[string]string{"foobar": ""},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			res := ParseAnnotations(test.Input)
			if diff := cmp.Diff(test.Expected, res); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
