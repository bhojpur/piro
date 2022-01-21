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
	"os"

	rice "github.com/GeertJohan/go.rice"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/godoc_vfs"
	log "github.com/sirupsen/logrus"
	"golang.org/x/tools/godoc/vfs/mapfs"
	"golang.org/x/xerrors"
)

// Migrate ensures that the database has the current schema required for using any of the postgres storage
func Migrate(db *sql.DB) error {
	fs, err := getMigrations(db)
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	mig, err := migrate.NewWithInstance("godoc-vfs", fs, "postgres", driver)
	if err != nil {
		return err
	}
	mig.Log = &logrusAdapter{}

	err = mig.Up()
	if err != nil && err != migrate.ErrNoChange {
		return xerrors.Errorf("error during migration: %w", err)
	}

	return nil
}

func getMigrations(db *sql.DB) (source.Driver, error) {
	box, err := rice.FindBox("migrations")
	if err != nil {
		return nil, err
	}
	migs := make(map[string]string)
	err = box.Walk("", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		migs[path], err = box.String(path)
		if err != nil {
			return xerrors.Errorf("cannot read from migration box: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, xerrors.Errorf("cannot list migrations: %w", err)
	}
	fs, err := godoc_vfs.WithInstance(mapfs.New(migs), "")
	if err != nil {
		return nil, err
	}

	return fs, nil
}

type logrusAdapter struct{}

func (*logrusAdapter) Printf(format string, args ...interface{}) {
	log.WithField("migration", true).Debugf(format, args...)
}

func (*logrusAdapter) Verbose() bool {
	return true
}
