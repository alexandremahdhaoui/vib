/*
Copyright 2023 Alexandre Mahdhaoui

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

package util

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alexandremahdhaoui/tooling/pkg/flaterrors"
	"github.com/alexandremahdhaoui/vib/internal/types"
)

// FileExist checks if a file exists.
func FileExist(path string) (bool, error) {
	err := MkBaseDir(path)
	if err != nil {
		return false, err
	}

	_, err = os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

// MkBaseDir creates the base directory for the given path.
func MkBaseDir(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
		return err
	}
	return nil
}

// Ptr returns a pointer to the given value.
func Ptr[T any](t T) *T {
	return &t
}

// RemoveIndexFromSlice removes an element from a slice at the given index.
func RemoveIndexFromSlice[T any](sl []T, i int) []T {
	sl[i] = sl[len(sl)-1]
	return sl[:len(sl)-1]
}

// JoinLine joins a line to a buffer.
func JoinLine(buffer string, line string) string {
	if line == "" {
		return buffer
	}

	if buffer == "" {
		return line
	}

	return strings.Join([]string{buffer, line}, "\n")
}

// EditFile opens a temporary file with the given content and allows the user to edit it with the given editor.
func EditFile(editor string, b []byte, encoding types.Encoding) ([]byte, error) {
	// we use encoding as filename extension for the user experience
	tmpFilename := fmt.Sprintf("vib-edit-*.%s", encoding)
	tmpfile, err := os.CreateTemp("", tmpFilename)
	if err != nil {
		return nil, flaterrors.Join(err, errors.New("creating temp file"))
	}
	defer os.Remove(tmpfile.Name()) // Clean up the temporary file later

	if _, err := tmpfile.Write(b); err != nil {
		return nil, flaterrors.Join(err, errors.New("writing to temp file"))
	}
	if err := tmpfile.Close(); err != nil {
		return nil, flaterrors.Join(err, errors.New("closing to temp file"))
	}

	cmd := exec.Command("sh", "-c", fmt.Sprintf("%s %s", editor, tmpfile.Name()))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf(
				"editor exited with error: %v, exit code: %d",
				exitError,
				exitError.ExitCode(),
			)
			// Handle specific exit codes if needed (e.g., user aborted without saving)
		} else {
			return nil, flaterrors.Join(errors.New("running editor"), err)
		}
	}

	out, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		return nil, flaterrors.Join(errors.New("Error reading edited file"), err)
	}

	return out, nil
}

// Must panics if the given error is not nil.
func Must[T any](out T, err error) T {
	if err != nil {
		panic(err.Error())
	}
	return out
}
