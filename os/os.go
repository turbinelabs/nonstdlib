/*
Copyright 2017 Turbine Labs, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package os provides an OS interface mirroring a subset of commonly used
// functions and variables from the golang os package.
package os

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE
//go:generate mockgen -destination mock_fileinfo.go -package $GOPACKAGE os FileInfo

import (
	"io"
	"io/ioutil"
	"os"
)

// OS is an interface on top of commonly used functions and variables
// from the golang os package, for easier testing. New methods can be
// added as needed.
//
// For undocumented functions, see correspondingly named things at
// https://golang.org/pkg/os.
type OS interface {
	Args() []string
	Getenv(key string) string
	LookupEnv(key string) (value string, found bool)
	ExpandEnv(s string) string
	Exit(code int)
	Stdin() io.Reader
	Stdout() io.Writer
	Stderr() io.Writer
	Stat(name string) (os.FileInfo, error)
	IsNotExist(err error) bool
	Remove(name string) error
	Rename(oldpath, newpath string) error
	Open(name string) (*os.File, error)
	OpenFile(name string, flag int, perm os.FileMode) (*os.File, error)

	// Constructs a new DirReader via NewDirReader.
	NewDirReader(dir string) DirReader
}

// New produces a new OS backed by the golang os package.
func New() OS {
	return goOS{}
}

type goOS struct{}

func (x goOS) Args() []string {
	return os.Args
}

func (x goOS) Getenv(key string) string {
	return os.Getenv(key)
}

func (x goOS) LookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}

func (x goOS) ExpandEnv(s string) string {
	return os.ExpandEnv(s)
}

func (x goOS) Exit(code int) {
	os.Exit(code)
}

func (x goOS) Stdin() io.Reader {
	return os.Stdin
}

func (x goOS) Stdout() io.Writer {
	return os.Stdout
}

func (x goOS) Stderr() io.Writer {
	return os.Stderr
}

func (x goOS) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (x goOS) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

func (x goOS) Remove(name string) error {
	return os.Remove(name)
}

func (x goOS) Rename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

func (x goOS) Open(name string) (*os.File, error) {
	return os.Open(name)
}

func (x goOS) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

func (x goOS) NewDirReader(dir string) DirReader {
	return NewDirReader(dir)
}

// ReadIfNonEmpty will read from an os.File only if non-empty; otherwise it will
// return an empty string. Useful, particularly, for checking for input on
// os.Stdin.
func ReadIfNonEmpty(f *os.File) (string, error) {
	st, err := f.Stat()
	if err != nil {
		return "", err
	}

	if st.Size() == 0 {
		return "", nil
	}

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
