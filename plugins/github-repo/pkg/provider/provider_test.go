package provider

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
			res := parseAnnotations(test.Input)
			if diff := cmp.Diff(test.Expected, res); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
