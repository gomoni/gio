# gio

Generic io - forked version of stdlib's io. Works with any slice of Go types.
Implements basic io interfaces and in-memory pipe and a unix-style pipeline.


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

 * split `gio` and `pipeline`?
 * figure out the error type, exit codes and so
 * add a unix specialization using type parameter `byte` making it compatible with
   Go standard library and unix expectations about a content of stdin/stdout
 * convert gonix and gonix/sbase on top of new core
 * explore the `Transform[F, T any]` pipe
