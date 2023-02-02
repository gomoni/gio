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

	err := NewPipeline[string]().Run(ctx, stdio, cat, wc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out.String())
	// Output: 3
```

## TODO

 * rework errors on top of go 1.20 ones
 * convert gonix and gonix/sbase on top of new core
 * os/exec helper for unix https://github.com/gomoni/gonix/blob/040661092859319d48d7664d99b1724eec64f636/pipe/exec.go
 * use some shlex helper for implement string support `cat | wc -l` would be
   executed by native code https://github.com/gomoni/gonix/blob/040661092859319d48d7664d99b1724eec64f636/pipe/sh.go
 * explore the `Transform[F, T any]` option allowing type conversion
