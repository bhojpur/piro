package prettyprint

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
	"text/tabwriter"
	"text/template"
	"time"

	tspb "google.golang.org/protobuf/types/known/timestamppb"
)

// TemplateFormat uses Go templates and tabwriter for formatting content
const TemplateFormat Format = "template"

func formatTemplate(pp *Content) error {
	tmpl, err := template.
		New("prettyprint").
		Funcs(map[string]interface{}{
			"toRFC3339": func(t *tspb.Timestamp) string {
				ts := tspb.Timestamp(*t)
				return ts.AsTime().Format(time.RFC3339)
			},
		}).
		Parse(pp.Template)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(pp.Writer, 8, 8, 8, ' ', 0)
	if err := tmpl.Execute(w, pp.Obj); err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return err
	}
	return nil
}
