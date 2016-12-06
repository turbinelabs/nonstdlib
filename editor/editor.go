/*
 * The editor package provides simple wrappers for interacting with an
 * environment configured text editor.
 *
 * The editor command is taken from the environment variable EDITOR and
 * the package variable DefaultEditor will be used if EDITOR is not set.
 */
package editor

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/turbinelabs/nonstdlib/os/exec"
)

// EditorVar specifies the environment variable we use to discover the editor
// command.
const EditorVar = "EDITOR"

var (
	// DefaultEditor is used if the user does not have EDITOR specified.
	DefaultEditor = ""

	// NoEditor indicates that no EDITOR was set in the environment.
	NoEditor = errors.New("could not find " + EditorVar + "environment variable")
)

func getEditor() (string, error) {
	edit := os.Getenv(EditorVar)
	if edit == "" && DefaultEditor != "" {
		edit = DefaultEditor
	}

	if edit == "" {
		return "", NoEditor
	}
	return edit, nil
}

// EditPath opens a specified path in the configured editor. If the path
// doesn't exist or fails to open then an error is returned.
func EditPath(path string) error {
	edit, err := getEditor()
	if err != nil {
		return err
	}

	f, err := os.Open(path)
	if f != nil {
		f.Close()
	}

	if err != nil {
		return err
	}

	return exec.RunInTerm(edit, path)
}

// EditText opens an editor with the initial contents populated from the
// provided string; the edited string is returned.
func EditText(str string) (string, error) {
	return EditTextType(str, "")
}

// EditTextType takes a string that should be edited and opens an editor with
// that content preloaded; it returns the edited string. If ext is not empty
// it will be used as the file extension.
//
// In the event of an error an empty string and the error will be returned.
func EditTextType(str, ext string) (string, error) {
	f, err := ioutil.TempFile("", "editor")
	if err != nil {
		return "", err
	}
	path := f.Name()
	f.Close()

	if ext != "" {
		newPath := fmt.Sprintf("%s.%s", path, ext)
		err := os.Rename(path, newPath)
		if err != nil {
			return "", err
		}
		path = newPath
	}

	defer os.Remove(path)

	err = ioutil.WriteFile(path, []byte(str), 0)
	if err != nil {
		return "", err
	}
	err = EditPath(path)

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
