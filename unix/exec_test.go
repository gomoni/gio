package unix_test

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/gomoni/gio/pipe"
	. "github.com/gomoni/gio/unix"
	"github.com/stretchr/testify/require"
)

func TestExec(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(func() {
		cancel()
	})

	stdout := bytes.NewBuffer(nil)
	stdio := NewStdio(
		nil,
		stdout,
		os.Stderr,
	)

	cmd := NewCmd(exec.Command("go", "version"))
	err := cmd.Run(ctx, stdio)
	require.NoError(t, err)
	require.True(t, strings.HasPrefix(stdout.String(), "go version"))
}

func TestNotFound(t *testing.T) {
	var path = "go"
	for {
		_, err := exec.LookPath(path)
		if err != nil {
			break
		}
		path = path + "a"
	}

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(func() {
		cancel()
	})

	stdout := bytes.NewBuffer(nil)
	stdio := NewStdio(
		nil,
		stdout,
		os.Stderr,
	)

	cmd := NewCmd(exec.Command(path, "version"))
	err := cmd.Run(ctx, stdio)
	require.Error(t, err)
	var pipeErr pipe.Error
	ok := errors.As(err, &pipeErr)
	require.True(t, ok)
	require.Equal(t, pipe.NotFound, pipeErr.Code)
	require.True(t, strings.Contains(pipeErr.Error(), "executable file not found"))
}
