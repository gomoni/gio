package gio

import (
	"strings"
)

type nopCloseR[T any] struct {
	r Reader[T]
}

func (n nopCloseR[T]) Read(data []T) (int, error) {
	return n.r.Read(data)
}

func (nopCloseR[T]) Close() error {
	return nil
}

type nopCloseW[T any] struct {
	w Writer[T]
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
	return s.errs[len(s.errs)-1]
}

func (s errorSlice) Error() string {
	ret := make([]string, 0, len(s.errs))
	for _, e := range s.errs {
		if e == nil {
			continue
		}
		ret = append(ret, e.Error())
	}
	return strings.Join(ret, "\n")
}
