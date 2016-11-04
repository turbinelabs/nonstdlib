package os

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"testing"

	"github.com/turbinelabs/test/assert"
)

type dirEntryMap map[string]interface{}

func mkFiles(t *testing.T, dir string, entries dirEntryMap) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err.Error())
	}
	for file, v := range entries {
		pathname := path.Join(dir, file)
		switch x := v.(type) {
		case dirEntryMap:
			mkFiles(t, pathname, x)
		case string:
			f, err := os.Create(pathname)
			if err == nil {
				_, err = f.WriteString(x)
				f.Close()
			}
			if err != nil {
				t.Fatal(err.Error())
			}
		default:
			t.Fatalf("invalid file entry type: %T, %s", v, v)
		}
	}
}

func prepDir(t *testing.T, entries dirEntryMap) (string, func()) {
	dir, err := ioutil.TempDir("", "dir-reader")
	if err != nil {
		t.Fatal(err.Error())
	}

	mkFiles(t, dir, entries)

	return dir, func() { os.RemoveAll(dir) }
}

func TestDirReaderRead(t *testing.T) {
	dirEntries := dirEntryMap{
		"a": "data",
		"b": "data",
		"c": "data",
		"d": dirEntryMap{"sub": "data"},
	}

	dir, cleanup := prepDir(t, dirEntries)
	defer cleanup()

	dirReader := NewDirReader(dir)
	files, err := dirReader.Read()
	assert.Nil(t, err)
	assert.HasSameElements(t, files, []string{"a", "b", "c", "d"})
}

func TestDirReaderReadEmpty(t *testing.T) {
	dir, cleanup := prepDir(t, dirEntryMap{})
	defer cleanup()

	dirReader := NewDirReader(dir)
	files, err := dirReader.Read()
	assert.Nil(t, err)
	assert.Equal(t, len(files), 0)
}

func TestDirReaderReadMany(t *testing.T) {
	dirEntries := dirEntryMap{}
	names := make([]string, 100)
	for i := range names {
		d := fmt.Sprintf("d%d", i)
		names[i] = d
		dirEntries[d] = "x"
	}

	dir, cleanup := prepDir(t, dirEntries)
	defer cleanup()

	dirReader := NewDirReader(dir)
	files, err := dirReader.Read()
	assert.Nil(t, err)
	assert.HasSameElements(t, files, names)
}

func TestDirReaderFilter(t *testing.T) {
	dirEntries := dirEntryMap{
		"a": "data",
		"b": "data",
		"c": dirEntryMap{"sub": "data"},
		"d": "data",
	}

	dir, cleanup := prepDir(t, dirEntries)
	defer cleanup()

	dirReader := NewDirReader(dir)
	files, err := dirReader.Filter(func(fi os.FileInfo) bool {
		return !fi.IsDir() && fi.Name() != "b"
	})
	assert.Nil(t, err)
	assert.HasSameElements(t, files, []string{"a", "d"})
}

func TestDirReaderFilterMany(t *testing.T) {
	dirEntries := dirEntryMap{}
	names := make([]string, 100)
	evenNames := make([]string, 0, 50)
	for i := range names {
		d := fmt.Sprintf("d%d", i)
		names[i] = d
		dirEntries[d] = "x"
		if i%2 == 0 {
			evenNames = append(evenNames, d)
		}
	}

	dir, cleanup := prepDir(t, dirEntries)
	defer cleanup()

	dirReader := NewDirReader(dir)
	files, err := dirReader.Filter(func(fi os.FileInfo) bool {
		i, err := strconv.Atoi(fi.Name()[1:])
		if err != nil {
			return false
		}
		return i%2 == 0
	})
	assert.Nil(t, err)
	assert.HasSameElements(t, files, evenNames)
}

func TestDirReaderError(t *testing.T) {
	dir, cleanup := prepDir(t, dirEntryMap{})
	cleanup()

	dirReader := NewDirReader(dir)
	files, err := dirReader.Read()
	assert.Nil(t, files)
	assert.NonNil(t, err)
}
