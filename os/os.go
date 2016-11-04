package os

//go:generate mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE
//go:generate mockgen -destination mock_fileinfo.go -package $GOPACKAGE os FileInfo

import (
	"io"
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

func New() OS {
	return goOS{}
}

type goOS struct{}

func (_ goOS) Args() []string {
	return os.Args
}

func (_ goOS) Getenv(key string) string {
	return os.Getenv(key)
}

func (_ goOS) LookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}

func (_ goOS) Exit(code int) {
	os.Exit(code)
}

func (_ goOS) Stdin() io.Reader {
	return os.Stdin
}

func (_ goOS) Stdout() io.Writer {
	return os.Stdout
}

func (_ goOS) Stderr() io.Writer {
	return os.Stderr
}

func (_ goOS) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (_ goOS) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

func (_ goOS) Remove(name string) error {
	return os.Remove(name)
}

func (_ goOS) Rename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

func (_ goOS) Open(name string) (*os.File, error) {
	return os.Open(name)
}

func (_ goOS) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

func (_ goOS) NewDirReader(dir string) DirReader {
	return NewDirReader(dir)
}
