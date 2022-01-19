package store_test

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
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bhojpur/piro/pkg/store"
)

func TestContinuousWriteReading(t *testing.T) {
	base, err := ioutil.TempDir(os.TempDir(), "tcwr")
	if err != nil {
		t.Errorf("cannot create test folder: %v", err)
	}

	s, err := store.NewFileLogStore(base)
	if err != nil {
		t.Errorf("cannot create test store: %v", err)
	}

	w, err := s.Open("foo")
	if err != nil {
		t.Errorf("cannot place log: %v", err)
	}
	r, err := s.Read("foo")
	if err != nil {
		t.Errorf("cannot read log: %v", err)
	}

	var msg = `hello world
	this is a test
	we're just writing stuff
	line by line`

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()

		lines := strings.Split(msg, "\n")
		for i, l := range lines {
			if i < len(lines)-1 {
				l += "\n"
			}

			n, err := w.Write([]byte(l))
			if err != nil {
				panic(fmt.Errorf("write error: %v", err))
			}
			if n != len(l) {
				panic(fmt.Errorf("write error: %v", io.ErrShortWrite))
			}
			time.Sleep(10 * time.Millisecond)
		}
		w.Close()
	}()

	rbuf := bytes.NewBuffer(nil)
	go func() {
		defer wg.Done()

		_, err := io.Copy(rbuf, r)
		if err != nil {
			t.Errorf("cannot read log: %+v", err)
			return
		}
	}()

	go func() {
		time.Sleep(5 * time.Second)
		panic("timeout")
	}()
	wg.Wait()

	actual := rbuf.Bytes()
	expected := []byte(msg)
	if !bytes.Equal(actual, expected) {
		for i, c := range actual {
			if i >= len(expected) {
				t.Errorf("read more than was written at byte %d: %v", i, c)
				continue
			}
			if c != expected[i] {
				t.Errorf("read difference at byte %d: %v !== %v", i, c, expected[i])
			}
		}
		t.Errorf("did not read message back, but: %s", string(actual))
	}
}
