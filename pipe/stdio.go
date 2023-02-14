package pipe

import (
	"context"
	"io"

	"github.com/gomoni/gio"
)

// pipe.StandardIO is a standard unix-like type safe input and output
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
