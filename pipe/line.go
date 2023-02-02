// Copyright 2023 Michal Vyskocil. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// package pipe implements type-safe unix like resources standard streams and pipelines
// This can be used to write unix like tools communicating through Go native types
// which can be connected together like unix allows.
package pipe

import (
	"context"
	"io"
	"sync"
	"sync/atomic"

	"github.com/gomoni/gio"
)

// gio.StandardIO is a standard unix-like type safe input and output
// defines three streams
//
//	stdin - from which can be read
//	stdout - to which results should be written
//	stderr - standard io.Writer for errors, debugs and so
type StandardIO[T any] interface {
	Stdin() gio.Reader[T]
	Stdout() gio.Writer[T]
	Stderr() io.Writer
}

type Filter[T any] interface {
	Run(context.Context, StandardIO[T]) error
}

// Stdio represent type safe unix-like standard input and output
// Implements gio.Standard[T] interface
type Stdio[T any] struct {
	stdin  gio.Reader[T]
	stdout gio.Writer[T]
	stderr io.Writer
}

func NewStdio[T any](stdin gio.Reader[T], stdout gio.Writer[T], stderr io.Writer) Stdio[T] {
	return Stdio[T]{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}
}

func (s Stdio[T]) Stdin() gio.Reader[T] {
	return s.stdin
}

func (s Stdio[T]) Stdout() gio.Writer[T] {
	return s.stdout
}

func (s Stdio[T]) Stderr() io.Writer {
	return s.stderr
}

// pipe.Line implements the unix-like pipeline for command | command2 | command3
// by connecting the filters via gio.Pipe
type Line[T any] struct {
	noPipeFail bool
}

func NewLine[T any]() Line[T] {
	return Line[T]{}
}

// Pipefail - true (the default) is an equivalent of set -o pipefail, so pipe is canceled
// on a first error and this is returned upper.
// To simulate default shell behavior use Pipefail(false). Which does not cancel the pipe and returns the
// last error.
func (p Line[T]) Pipefail(b bool) Line[T] {
	p.noPipeFail = !b
	return p
}

// Run joins all filters via gio.Pipe with a stdio. Each filter runs in own goroutine and function
// returns on all. The returned error depends on Pipefail value
//
// true (the default) - returns nil if none of commands fail, otherwise returns
// all errors in a pipe in a slice. If the first failure is Error, then it's
// Code is returned. otherwise code 1 is used.
//
// false - returns error of a last command. If the error is nil, then result of
// the call is nil. for non nil errors, it returns a slice of all errors and a
// code or 1 depending on a type of last error.
func (p Line[T]) Run(ctx context.Context, stdio StandardIO[T], filters ...Filter[T]) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if len(filters) == 1 {
		return filters[0].Run(ctx, stdio)
	}

	errs := errorSlice{errs: make([]error, len(filters))}
	var hasError atomic.Bool
	var wg sync.WaitGroup
	var in gio.ReadCloser[T] = nopCloseR[T]{r: stdio.Stdin()}
	for idx, filter := range filters {
		var nextIn gio.ReadCloser[T]
		var out gio.WriteCloser[T]
		isLast := idx == len(filters)-1
		if isLast {
			out = nopCloseW[T]{w: stdio.Stdout()}
		} else {
			pipeR, pipeW := gio.Pipe[T]()
			out = pipeW
			nextIn = pipeR
		}

		wg.Add(1)
		go p.runOne(
			ctx,
			cancel,
			&errs,
			&hasError,
			idx,
			&wg,
			filter,
			gostdio[T]{stdin: in, stdout: out, stderr: stdio.Stderr()})
		in = nextIn
	}

	wg.Wait()

	// XXX: improve error handling
	//      wait on 1.20 errors with Join and Unwrap []error
	if p.noPipeFail {
		return errs.noPipefail(1)
	} else {
		return errs.pipefail(1)
	}
}

type gostdio[T any] struct {
	stdin  gio.ReadCloser[T]
	stdout gio.WriteCloser[T]
	stderr io.Writer
}

func (p Line[T]) runOne(ctx context.Context, cancel context.CancelFunc, errs *errorSlice, hasError *atomic.Bool, idx int, wg *sync.WaitGroup, filter Filter[T], stdio gostdio[T]) {
	defer wg.Done()
	defer stdio.stdin.Close()
	defer stdio.stdout.Close()

	// do not start more tasks
	if !p.noPipeFail && hasError.Load() {
		return
	}

	err := filter.Run(ctx, Stdio[T]{stdin: stdio.stdin, stdout: stdio.stdout, stderr: stdio.stderr})
	errs.set(idx, err)
	if err != nil {
		hasError.Store(true)
		if !p.noPipeFail {
			cancel()
		}
	}
}
