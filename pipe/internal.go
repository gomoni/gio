// Copyright 2023 Michal Vyskocil. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package pipe

import (
	"errors"

	"github.com/gomoni/gio"
)

type nopCloseR[T any] struct {
	r gio.Reader[T]
}

func (n nopCloseR[T]) Read(data []T) (int, error) {
	return n.r.Read(data)
}

func (nopCloseR[T]) Close() error {
	return nil
}

type nopCloseW[T any] struct {
	w gio.Writer[T]
}

func (n nopCloseW[T]) Write(data []T) (int, error) {
	return n.w.Write(data)
}

func (nopCloseW[T]) Close() error {
	return nil
}

type errorSlice struct {
	errs []error
}

func (s *errorSlice) set(idx int, err error) {
	s.errs[idx] = err
}

func (s errorSlice) Unwrap() []error {
	return s.errs
}

func (s errorSlice) firstNonNil() error {
	for _, err := range s.errs {
		if err == nil {
			continue
		}
		return err
	}
	return nil
}

func (s errorSlice) last() error {
	if len(s.errs) == 0 {
		return nil
	}
	return s.errs[len(s.errs)-1]
}

func (s errorSlice) Error() string {
	return errors.Join(s.errs...).Error()
}

func (s errorSlice) noPipefail(code int) error {
	err := s.last()
	if err == nil {
		return nil
	}
	var pipeError Error
	if errors.As(err, &pipeError) {
		pipeError.Err = s
		return pipeError
	}
	return NewError(code, s)
}

func (s errorSlice) pipefail(code int) error {
	err := s.firstNonNil()
	if err == nil {
		return nil
	}
	var pipeError Error
	if !errors.As(err, &pipeError) {
		pipeError = NewError(code, s)
	} else {
		pipeError.Err = s
	}

	return pipeError
}
