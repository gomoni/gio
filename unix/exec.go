package unix

import (
	"context"
	"os/exec"

	"github.com/gomoni/gio/pipe"
)

// Exec is a wrapper of os/exec.Cmd providing a Filter compatible interface
type Cmd struct {
	cmd *exec.Cmd
}

// NewCmd wraps exec.Cmd. The intended usage is
//
//	cmd := NewCmd(exec.Command("go", "version"))
func NewCmd(cmd *exec.Cmd) Cmd {
	if cmd == nil {
		panic("cmd is nil")
	}
	return Cmd{cmd}
}

// Run implements Filter interface for Cmd wrapper. It creates a _new_ instance of
// exec.Command under a hood, so any previous state is ignored here.
//
// Returns a [pipe.Error] if Run results in [*exec.ExitError]. Code is ExitCode and
// Err is the *exec.ExitError
func (c Cmd) Run(ctx context.Context, stdio StandardIO) error {
	cmd := exec.CommandContext(ctx, c.cmd.Path, c.cmd.Args[1:]...)
	cmd.Env = c.cmd.Env
	cmd.Dir = c.cmd.Dir

	cmd.Stdin = stdio.Stdin()
	cmd.Stdout = stdio.Stdout()
	cmd.Stderr = stdio.Stderr()

	cmd.ExtraFiles = c.cmd.ExtraFiles
	cmd.SysProcAttr = c.cmd.SysProcAttr
	cmd.WaitDelay = c.cmd.WaitDelay

	err := cmd.Run()
	if err == nil {
		return nil
	}
	return pipe.FromError(err)
}
