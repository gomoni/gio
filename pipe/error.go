// Copyright 2022 Michal Vyskocil. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package pipe

import (
	"errors"
	"fmt"
	"io/fs"
	"os/exec"
)

const (
	// NotExecutable is from POSIX and indicates tool was found, but not executable
	NotExecutable = 126
	// NotFound is from POSIX and indicate a tool was not found
	NotFound = 127
	// UnknownError is a code used for unpacking other than pipe.Error
	UnknownError = 250
)

// Error is a common error type returned by pipeline. It has a Code for unix
// compatibility and an error.
type Error struct {
	Code int
	Err  error
}

func (e Error) Error() string {
	return fmt.Sprintf("Error{Code: %d, Err: %+v}", e.Code, e.Err)
}

func (e Error) Unwrap() error {
	return e.Err
}

// Errors returns a slice of errors if err is Error and member Err implements
// Unwrap []error otherwise returns nil even for non-nil errors
func Errors(err error) []error {

	if Err, ok := err.(Error); ok {
		if errs, ok := Err.Err.(interface{ Unwrap() []error }); ok {
			return errs.Unwrap()
		}
	}
	return nil
}

// NewError returns a new error with code and error inside
func NewError(code int, err error) Error {
	return Error{Code: code, Err: err}
}

// NewErrorf returns formatted new error with code and error inside
func NewErrorf(code int, format string, args ...any) Error {
	return Error{Code: code, Err: fmt.Errorf(format, args...)}
}

// FromError unpacks error into Error. If it can't be unpacked, it assigns code 250
// error fs.ErrPermission will get code NotExecutable (126)
// error exec.ErrNotFound will get code NotFound (code 127)
func FromError(x error) Error {
	if errors.Is(x, exec.ErrNotFound) {
		return Error{Code: NotFound, Err: x}
	}

	var fsErr *fs.PathError
	if errors.As(x, &fsErr) {
		if fsErr.Op == "fork/exec" && errors.Is(fsErr.Err, fs.ErrPermission) {
			return Error{Code: NotExecutable, Err: x}
		}
	}

	var err Error
	if !errors.As(x, &err) {
		return Error{Code: UnknownError, Err: x}
	}
	return err
}
