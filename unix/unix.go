// Copyright 2023 Michal Vyskocil. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// package pipe implements unix like resources standard streams and pipelines
// This can be used to write unix like tools communicating through Go stream of bytes
// which can be connected together like unix allows.
//
// It is build as a tiny wrapper on top of gio/pipe as it uses standard io.Reader and io.Writer
// so it can integrate into idiomatic Go code more seamlessly.
package unix

import (
	"context"
	"io"

	"github.com/gomoni/gio"
	"github.com/gomoni/gio/pipe"
)

type StandardIO interface {
	Stdin() io.Reader
	Stdout() io.Writer
	Stderr() io.Writer
}

type Filter interface {
	Run(context.Context, StandardIO) error
}

type Stdio struct {
	stdin  gio.Reader[byte]
	stdout gio.Writer[byte]
	stderr gio.Writer[byte]
}

func NewStdio(stdin io.Reader, stdout io.Writer, stderr io.Writer) Stdio {
	return Stdio{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}
}

func (s Stdio) Stdin() io.Reader {
	return s.stdin
}

func (s Stdio) Stdout() io.Writer {
	return s.stdout
}

func (s Stdio) Stderr() io.Writer {
	return s.stderr
}

type Line struct {
	pipe.Line[byte]
}

func NewLine() Line {
	return Line{Line: pipe.NewLine[byte]()}
}

func (p Line) Pipefail(b bool) Line {
	return p.Pipefail(b)
}

func (p Line) Run(ctx context.Context, stdio StandardIO, filters ...Filter) error {
	pipeio := pipe.NewStdio[byte](stdio.Stdin(), stdio.Stdout(), stdio.Stderr())
	pipefilters := make([]pipe.Filter[byte], len(filters))
	for idx, f := range filters {
		pipefilters[idx] = pipeFilter{filter: f}
	}

	return p.Line.Run(ctx, pipeio, pipefilters...)
}

type pipeFilter struct {
	filter Filter
}

func (f pipeFilter) Run(ctx context.Context, stdio pipe.StandardIO[byte]) error {
	unixio := NewStdio(stdio.Stdin(), stdio.Stdout(), stdio.Stderr())
	return f.filter.Run(ctx, unixio)
}
