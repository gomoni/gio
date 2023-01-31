package gio_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/gomoni/gio"
)

func mustWrite[T any](t *testing.T, w Writer[T], data []T, doneCh chan struct{}) {
	t.Helper()
	n, err := w.Write(data)
	require.NoError(t, err)
	require.Equal(t, len(data), n)
	close(doneCh)
}

// test single write/read call
func TestPipe1(t *testing.T) {
	doneCh := make(chan struct{})
	rd, wr := Pipe[int]()
	t.Cleanup(func() {
		rd.Close()
		wr.Close()
	})

	go mustWrite[int](t, wr, []int{42}, doneCh)
	res := []int{0}
	n, err := rd.Read(res)
	require.NoError(t, err)
	require.Equal(t, 1, n)
	require.Equal(t, 42, res[0])
	<-doneCh
}
