package os

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE

import (
	"io"
	"os"
)

const numEntriesPerRead = 32

type DirEntryFilter func(os.FileInfo) bool

type DirReader interface {
	Read() ([]string, error)

	Filter(DirEntryFilter) ([]string, error)
}

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
