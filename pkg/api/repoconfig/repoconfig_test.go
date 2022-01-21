package repoconfig_test

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
	"encoding/json"
	"testing"

	"github.com/bhojpur/piro/pkg/api/repoconfig"
	v1 "github.com/bhojpur/piro/pkg/api/v1"
	"gopkg.in/yaml.v3"
)

func TestUnmarshalC(t *testing.T) {
	tests := []struct {
		Source      string
		Expectation string
	}{
		{`defaultJob: "foo.yaml"`, `{"DefaultJob":"foo.yaml","Rules":null}`},
		{
			`rules:
- path: ""
  matchesAll:
  - or: ["repo.ref ~= refs/tags/"]
- path: ""
  matchesAll:
  - or: ["repo.ref !~= refs/branches/"]`,
			`{"DefaultJob":"","Rules":[{"Path":"","Expr":[{"terms":[{"field":"repo.ref","value":"refs/tags/","operation":3}]}]},{"Path":"","Expr":[{"terms":[{"field":"repo.ref","value":"refs/branches/","operation":3,"negate":true}]}]}]}`,
		},
		{
			`rules:
- path: "foo.yaml"
  matchesAll:
  - or:
    - "repo.ref ~= refs/branches/"
  - or:
    - "name !~= 0"
`, `{"DefaultJob":"","Rules":[{"Path":"foo.yaml","Expr":[{"terms":[{"field":"repo.ref","value":"refs/branches/","operation":3}]},{"terms":[{"field":"name","value":"0","operation":3,"negate":true}]}]}]}`,
		},
	}

	for idx, test := range tests {
		var c repoconfig.C
		err := yaml.Unmarshal([]byte(test.Source), &c)
		if err != nil {
			t.Errorf("test %d: %v", idx, err)
			continue
		}

		act, err := json.Marshal(c)
		if err != nil {
			t.Errorf("test %d: %v", idx, err)
			continue
		}

		if string(act) != test.Expectation {
			t.Errorf("test %d: did not match expectation.\nExpected: %s\nActual: %s\n", idx, test.Expectation, string(act))
		}
	}
}

func TestTemplatePath(t *testing.T) {
	tests := []struct {
		C repoconfig.C
		M v1.JobMetadata
		E string
	}{
		{repoconfig.C{}, v1.JobMetadata{}, ""},
		{repoconfig.C{}, v1.JobMetadata{Owner: "foo", Repository: &v1.Repository{Owner: "foo"}, Trigger: v1.JobTrigger_TRIGGER_MANUAL}, ""},
		{repoconfig.C{DefaultJob: "foo"}, v1.JobMetadata{}, "foo"},
		{repoconfig.C{DefaultJob: "foo", Rules: []*repoconfig.JobStartRule{&repoconfig.JobStartRule{Path: "bar"}}}, v1.JobMetadata{}, "bar"},
		{
			repoconfig.C{
				DefaultJob: "foo",
				Rules: []*repoconfig.JobStartRule{
					&repoconfig.JobStartRule{
						Path: "bar",
						Expr: []*v1.FilterExpression{
							&v1.FilterExpression{Terms: []*v1.FilterTerm{&v1.FilterTerm{Field: "repo.ref", Value: "test", Operation: v1.FilterOp_OP_EQUALS}}},
						},
					},
				},
			},
			v1.JobMetadata{},
			"foo",
		},
		{
			repoconfig.C{
				DefaultJob: "foo",
				Rules: []*repoconfig.JobStartRule{
					&repoconfig.JobStartRule{
						Path: "bar",
						Expr: []*v1.FilterExpression{
							&v1.FilterExpression{Terms: []*v1.FilterTerm{&v1.FilterTerm{Field: "repo.ref", Value: "test", Operation: v1.FilterOp_OP_EQUALS}}},
						},
					},
				},
			},
			v1.JobMetadata{
				Repository: &v1.Repository{
					Ref: "test",
				},
			},
			"bar",
		},
	}

	for idx, test := range tests {
		act := test.C.TemplatePath(&test.M)
		if act != test.E {
			t.Errorf("test %d: expected %s, actual %s", idx, test.E, act)
		}
	}
}
