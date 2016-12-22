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
	"strings"

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

// Get returns the editor that will be used as determined by reading the
// EditorVar environment variable. The unmodified content of EditorVar will
// returned as the first return value. The second return value will be a
// space split slice of the contents.
func Get() (string, []string, error) {
	return getEditor()
}

func getEditor() (string, []string, error) {
	edit := os.Getenv(EditorVar)
	if edit == "" && DefaultEditor != "" {
		edit = DefaultEditor
	}
	if edit == "" {
		return "", nil, NoEditor
	}
	parts := strings.Split(edit, " ")

	return edit, parts, nil
}

// EditPath opens a specified path in the configured editor. If the path
// doesn't exist or fails to open then an error is returned.
func EditPath(path string) error {
	_, tokenized, err := getEditor()
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

	tokenized = append(tokenized, path)
	return exec.RunInTerm(tokenized[0], tokenized[1:]...)
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
