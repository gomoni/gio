moved to https://codeberg.org/gonix/gio
# gio

 * [gio](.): Generic-aware io interfaces like `gio.Reader[int]`, generic-aware in-memory pipe `gio.Pipe[string]`. Forked from Go stdlib.
 * [gio/pipe](./pipe): Generic-aware pipeline with a standard input output streams and filters. Enable writing unix-like utilities working on top of native Go types.
 * [gio/unix](./unix): byte stream aware pipeline with a standard input output streams and filters. Works like traditional unix tools.

## Example

An equivalent of `cat | wc -l` using a native Go types and channels under the hood.

```go
	out := &StringBuffer{}
	stdio := NewStdio[string](
		nil,
		out,
		os.Stderr,
	)

	err := pipe.NewLine[string]().Run(ctx, stdio, cat, wc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out.String())
	// Output: 3
```

## os/exec wrapper

`gio/unix` has a `*exec.Cmd` wrapper allowing to run any system command as a Filter

```
	stdout := bytes.NewBuffer(nil)
	stdio := unix.NewStdio(
		nil,
		stdout,
		os.Stderr,
	)

	cmd := unix.NewCmd(exec.Command("go", "version"))
	err := cmd.Run(ctx, stdio)
	fmt.Println(out.String())
	// Output: go version 1.20 linux/amd64
```


## TODO

 * explore the `Transform[F, T any]` option allowing type conversion
