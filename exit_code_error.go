/*
 Copyright 2024 The lvm2go Authors.

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

package lvm2go

import (
	"bytes"
	"errors"
	"fmt"
)

// AsExitCodeError returns the ExitCodeError from the error if it exists and a bool indicating if is an ExitCodeError or not.
func AsExitCodeError(err error) (ExitCodeError, bool) {
	var exitCodeErr ExitCodeError
	ok := errors.As(err, &exitCodeErr)
	return exitCodeErr, ok
}

// ExitCodeError is an error that wraps the original error and the stderr output of the command run.
// It also provides an exit code if present that can be used to determine the type of error from LVM.
// Regular inaccessible errors will have an exit code of 5.
type ExitCodeError interface {
	error
	ExitCode() int
	Unwrap() error
}

// NewExitCodeError returns a new ExitCodeError with the provided error and stderr output.
func NewExitCodeError(err error, stderr []byte) ExitCodeError {
	return &exitCodeErr{
		err:    err,
		stderr: stderr,
	}
}

// exitCodeErr is an implementation of ExitCodeError storing the original error and the stderr output of the lvmBinaryPath command.
// It also provides a POSIX exit code that can be used to determine the type of error from LVM.
type exitCodeErr struct {
	err    error
	stderr []byte
}

func (e *exitCodeErr) Error() string {
	if e.stderr != nil {
		return fmt.Sprintf("%v: %v", e.err, string(bytes.TrimSpace(e.stderr)))
	}
	return e.err.Error()
}

func (e *exitCodeErr) Unwrap() error {
	return e.err
}

func (e *exitCodeErr) ExitCode() int {
	type exitError interface {
		ExitCode() int
		error
	}
	var err exitError
	if errors.As(e.err, &err) {
		return err.ExitCode()
	}
	return -1
}
