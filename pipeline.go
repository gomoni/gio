package gio

import (
	"context"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
)

// gio.Standard is a standard unix-like type safe input and output
// defines three streams
//
//	stdin - from which can be read
//	stdout - to which results should be written
//	stderr - standard io.Writer for errors, debugs and so
type Standard[T any] interface {
	Stdin() Reader[T]
	Stdout() Writer[T]
	Stderr() io.Writer
}

type Filter[T any] interface {
	Run(context.Context, Standard[T]) error
}

type Pipeline[T any] struct {
	noPipeFail bool
}

func NewPipeline[T any]() Pipeline[T] {
	return Pipeline[T]{}
}

func (p Pipeline[T]) Pipefail(b bool) Pipeline[T] {
	p.noPipeFail = !b
	return p
}

func (p Pipeline[T]) Run(ctx context.Context, stdio Standard[T], filters ...Filter[T]) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if len(filters) == 1 {
		return filters[0].Run(ctx, stdio)
	}

	errs := errorSlice{errs: make([]error, len(filters))}
	var hasError atomic.Bool
	var wg sync.WaitGroup
	var in ReadCloser[T] = nopCloseR[T]{r: stdio.Stdin()}
	for idx, filter := range filters {
		var nextIn ReadCloser[T]
		var out WriteCloser[T]
		isLast := idx == len(filters)-1
		if isLast {
			out = nopCloseW[T]{w: stdio.Stdout()}
		} else {
			pipeR, pipeW := Pipe[T]()
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
	//      port the error with exit code
	//      wait on 1.20 errors with Join and Unwrap []error
	//      define what to do for noPipeFail - but non last command returns an error
	var err error
	if p.noPipeFail {
		err = errs.firstNonNil()
	} else {
		err = errs.last()
	}
	if hasError.Load() {
		if err == nil {
			return &errs
		}
		return fmt.Errorf("%w: %s", err, errs.Error())
	}
	return nil
}

type gostdio[T any] struct {
	stdin  ReadCloser[T]
	stdout WriteCloser[T]
	stderr io.Writer
}

func (p Pipeline[T]) runOne(ctx context.Context, cancel context.CancelFunc, errs *errorSlice, hasError *atomic.Bool, idx int, wg *sync.WaitGroup, filter Filter[T], stdio gostdio[T]) {
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
		fmt.Fprintf(stdio.stderr, "%+v", err)
	}
}
