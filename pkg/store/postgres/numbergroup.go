package postgres

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
	"database/sql"

	"github.com/bhojpur/piro/pkg/store"
)

// NumberGroup provides postgres backed number groups
type NumberGroup struct {
	DB *sql.DB
}

// NewNumberGroup creates a new SQL number group store
func NewNumberGroup(db *sql.DB) (*NumberGroup, error) {
	return &NumberGroup{DB: db}, nil
}

// Latest returns the latest number of a particular number group.
func (ngrp *NumberGroup) Latest(group string) (nr int, err error) {
	err = ngrp.DB.QueryRow(`
		SELECT val
		FROM   number_group
		WHERE  name = $1`,
		group,
	).Scan(&nr)
	if err == sql.ErrNoRows {
		return 0, store.ErrNotFound
	}
	return
}

// Next returns the next number in the group.
func (ngrp *NumberGroup) Next(group string) (nr int, err error) {
	err = ngrp.DB.QueryRow(`
		INSERT
		INTO   number_group (name, val)
		VALUES              ($1  , 0  )
		ON CONFLICT (name) DO UPDATE 
			SET val = number_group.val + 1
		RETURNING val`,
		group,
	).Scan(&nr)
	return
}
