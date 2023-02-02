// Copyright 2023 Michal Vyskocil. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package unix_test

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	. "github.com/gomoni/gio/unix"
)

type Cat struct {
	cat []byte
}

func (c Cat) Run(ctx context.Context, stdio StandardIO) error {
	buf := bytes.NewBuffer(c.cat)
	_, err := io.Copy(stdio.Stdout(), buf)
	return err
}

// CountLines count a number of read lines
type CountLines struct{}

func (c CountLines) Run(ctx context.Context, stdio StandardIO) error {
	counter := 0
	r := bufio.NewScanner(stdio.Stdin())
	for r.Scan() {
		counter++
	}
	if r.Err() != nil {
		return r.Err()
	}
	fmt.Fprintf(stdio.Stdout(), "%d\n", counter)
	return nil
}

type Fail struct {
	err error
}

func (c Fail) Run(context.Context, StandardIO) error {
	return c.err
}

func Example() {
	ctx := context.Background()
	cat := Cat{
		cat: []byte("three\nsmall\npigs\n"),
	}
	wc := CountLines{}

	var out strings.Builder
	stdio := NewStdio(
		nil,
		&out,
		os.Stderr,
	)

	// an equivalent of cat | wc -l
	err := NewLine().Run(ctx, stdio, cat, wc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out.String())
	// Output: 3
}
