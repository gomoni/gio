// Copyright 2023 Michal Vyskocil. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package pipe_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	. "github.com/gomoni/gio/pipe"
)

// Lines sends each line to stdout
type Lines struct {
	cat []string
}

func (c Lines) Run(ctx context.Context, stdio StandardIO[string]) error {
	for _, line := range c.cat {
		_, err := stdio.Stdout().Write([]string{line})
		if err != nil {
			return err
		}
	}
	return nil
}

// CountLines count a number of read lines
type CountLines struct {
}

func (c CountLines) Run(ctx context.Context, stdio StandardIO[string]) error {
	counter := 0
	for {
		var s []string = []string{""}
		_, err := stdio.Stdin().Read(s)
		if errors.Is(err, io.EOF) {
			_, err := stdio.Stdout().Write([]string{strconv.Itoa(counter)})
			return err
		} else if err != nil {
			return err
		}
		counter++
	}
}

type Fail struct {
	err error
}

func (c Fail) Run(context.Context, StandardIO[string]) error {
	return c.err
}

// StringBuffer implements the Writer[string] interface
type StringBuffer struct {
	s strings.Builder
}

func (s *StringBuffer) Write(str []string) (int, error) {
	for idx, str := range str {
		_, err := s.s.WriteString(str)
		if err != nil {
			return idx + 1, err
		}
	}
	return len(str), nil
}

func (s StringBuffer) String() string {
	return s.s.String()
}

func Example() {
	ctx := context.Background()
	cat := Lines{
		cat: []string{"three", "small", "pigs"},
	}
	wc := CountLines{}

	out := &StringBuffer{}
	stdio := NewStdio[string](
		nil,
		out,
		os.Stderr,
	)

	// an equivalent of cat | wc -l
	// just using a native Go types and channels
	err := NewLine[string]().Run(ctx, stdio, cat, wc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out.String())
	// Output: 3
}

func ExampleLine_pipefail() {
	ctx := context.Background()
	cat := Lines{
		cat: []string{"three", "small", "pigs"},
	}
	fail := Fail{err: io.EOF}
	wc := CountLines{}

	out := &StringBuffer{}
	stdio := NewStdio[string](
		nil,
		out,
		os.Stderr,
	)

	// an equivalent of set -o pipefail; false | cat | wc -l
	err := NewLine[string]().Pipefail(true).Run(ctx, stdio, fail, cat, wc)
	if err == nil {
		log.Fatal("expected err, got nil")
	}
	fmt.Println(err)
	// Output: Error{Code: 1, Err: EOF}
}

func ExampleLine_nopipefail() {
	ctx := context.Background()
	cat := Lines{
		cat: []string{"three", "small", "pigs"},
	}
	fail := Fail{err: io.EOF}
	wc := CountLines{}

	out := &StringBuffer{}
	stdio := NewStdio[string](
		nil,
		out,
		os.Stderr,
	)

	// an equivalent of false | cat | wc -l
	err := NewLine[string]().Pipefail(false).Run(ctx, stdio, fail, cat, wc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("OK")
	// Output: OK
}
