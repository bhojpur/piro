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
	"fmt"
	"io"

	"google.golang.org/protobuf/proto"
)

// Format defines the kind of pretty-printing format we want to use
type Format string

// HasFormat returns true if the format is supported
func HasFormat(fmt Format) bool {
	_, ok := formatter[fmt]
	return ok
}

const (
	// StringFormat uses the Go-builtin stringification for printing
	StringFormat Format = "string"
)

type formatterFunc func(*Content) error

var formatter = map[Format]formatterFunc{
	StringFormat:   formatString,
	TemplateFormat: formatTemplate,
	JSONFormat:     formatJSON,
	YAMLFormat:     formatYAML,
}

func formatString(pp *Content) error {
	_, err := fmt.Fprintf(pp.Writer, "%s", pp.Obj)
	return err
}

// Content is pretty-printable content
type Content struct {
	Obj      proto.Message
	Format   Format
	Writer   io.Writer
	Template string
}

// Print outputs the content to its writer in the given format
func (pp *Content) Print() error {
	formatter, ok := formatter[pp.Format]
	if !ok {
		return fmt.Errorf("Unknown format: %s", pp.Format)
	}

	return formatter(pp)
}
