/*
Copyright 2018 Turbine Labs, Inc.

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

package os

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE --write_package_comment=false

import (
	"io"
	"os"
)

const numEntriesPerRead = 32

// DirEntryFilter is a function that filters directory entries by
// returning true for directories that should be accepted.
type DirEntryFilter func(os.FileInfo) bool

// DirReader reads or reads and filters directory entries.
type DirReader interface {
	// Read reads all entries from the directory.
	Read() ([]string, error)

	// Filter returns all entries from the directory that match
	// the given filter.
	Filter(DirEntryFilter) ([]string, error)
}

// NewDirReader constructs a new DirReader for the given path name.
func NewDirReader(pathname string) DirReader {
	return &dirReader{pathname}
}

type dirReader struct {
	pathname string
}

func (dr *dirReader) Read() ([]string, error) {
	return dr.Filter(func(_ os.FileInfo) bool { return true })
}

func (dr *dirReader) Filter(accept DirEntryFilter) ([]string, error) {
	dirHandle, err := os.Open(dr.pathname)
	if err != nil {
		return nil, err
	}
	defer dirHandle.Close()

	result := []string{}
	var fileInfos []os.FileInfo
	err = nil
	for err == nil {
		fileInfos, err = dirHandle.Readdir(numEntriesPerRead)
		for _, fileInfo := range fileInfos {
			if accept(fileInfo) {
				result = append(result, fileInfo.Name())
			}
		}
	}
	if err == io.EOF {
		err = nil
	}

	return result, err
}
