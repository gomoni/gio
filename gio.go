// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copyright 2023 Michal Vyskocil. All rights reserved.

// Package gio - generic io - is forked variant of Go stdlib io module.
// Code is extended to be able to work on a slice of Go any type.
package gio

// Reader is the interface that wraps the basic Read method.
//
// Read reads up to len(p) types into p. It returns the number of types
// read (0 <= n <= len(p)) and any error encountered. Even if Read
// returns n < len(p), it may use all of p as scratch space during the call.
// If some data is available but not len(p) bytes, Read conventionally
// returns what is available instead of waiting for more.
//
// When Read encounters an error or end-of-file condition after
// successfully reading n > 0 types, it returns the number of
// bytes read. It may return the (non-nil) error from the same call
// or return the error (and n == 0) from a subsequent call.
// An instance of this general case is that a Reader returning
// a non-zero number of types at the end of the input stream may
// return either err == EOF or err == nil. The next Read should
// return 0, EOF.
//
// Callers should always process the n > 0 bytes returned before
// considering the error err. Doing so correctly handles I/O errors
// that happen after reading some types and also both of the
// allowed EOF behaviors.
//
// Implementations of Read are discouraged from returning a
// zero byte count with a nil error, except when len(p) == 0.
// Callers should treat a return of 0 and nil as indicating that
// nothing happened; in particular it does not indicate EOF.
//
// Implementations must not retain p.
type Reader[T any] interface {
	Read(p []T) (n int, err error)
}

// Writer is the interface that wraps the basic Write method.
//
// Write writes len(p) types from p to the underlying data stream.
// It returns the number of types written from p (0 <= n <= len(p))
// and any error encountered that caused the write to stop early.
// Write must return a non-nil error if it returns n < len(p).
// Write must not modify the slice data, even temporarily.
//
// Implementations must not retain p.
type Writer[T any] interface {
	Write(p []T) (n int, err error)
}

// Closer is the interface that wraps the basic Close method.
//
// The behavior of Close after the first call is undefined.
// Specific implementations may document their own behavior.
type Closer interface {
	Close() error
}

// ReadCloser is the interface that groups the basic Read and Close methods.
type ReadCloser[T any] interface {
	Reader[T]
	Closer
}

// WriteCloser is the interface that groups the basic Write and Close methods.
type WriteCloser[T any] interface {
	Writer[T]
	Closer
}
