package gio

import "io"

// Stdio represent type safe unix-like standard input and output
// Implements gio.Standard[T] interface
type Stdio[T any] struct {
	stdin  Reader[T]
	stdout Writer[T]
	stderr io.Writer
}

func NewStdio[T any](stdin Reader[T], stdout Writer[T], stderr io.Writer) Stdio[T] {
	return Stdio[T]{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}
}

func (s Stdio[T]) Stdin() Reader[T] {
	return s.stdin
}

func (s Stdio[T]) Stdout() Writer[T] {
	return s.stdout
}

func (s Stdio[T]) Stderr() io.Writer {
	return s.stderr
}
