package process

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type RunningProcess struct {
	cmd          *exec.Cmd
	stdout       bytes.Buffer
	stderr       bytes.Buffer
	stdOutReader io.ReadCloser
	stdErrReader io.ReadCloser
}

func (p *RunningProcess) Wait() error {
	return p.cmd.Wait()
}

func (p *RunningProcess) Stdout() string {
	return p.stdout.String()
}

func (p *RunningProcess) Stderr() string {
	return p.stderr.String()
}

func (p *RunningProcess) Kill() error {
	p.stdErrReader.Close()
	p.stdOutReader.Close()
	return p.cmd.Process.Signal(os.Interrupt)
}

type CmdOpt func(*exec.Cmd)

func WithDir(dir string) CmdOpt {
	return func(cmd *exec.Cmd) {
		cmd.Dir = dir
	}
}

func ExecuteCommand(ctx context.Context, cmdName string, args []string, opts ...CmdOpt) (*RunningProcess, error) {
	cmd := exec.CommandContext(ctx, cmdName, args...)
	for _, opt := range opts {
		opt(cmd)
	}

	stdOutReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error creating StdoutPipe for Cmd", err)
		os.Exit(1)
	}
	stdErrReader, err := cmd.StderrPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error creating StderrPipe for Cmd", err)
		os.Exit(1)
	}

	var outb, errb bytes.Buffer

	stdOutScanner := bufio.NewScanner(stdOutReader)
	go func() {
		for stdOutScanner.Scan() {
			t := stdOutScanner.Text()
			if len(t) == 0 {
				continue
			}
			outb.WriteString(t)
			outb.WriteByte('\n')
			fmt.Fprintln(os.Stdout, t)
		}
	}()

	stdErrScanner := bufio.NewScanner(stdErrReader)
	go func() {
		for stdErrScanner.Scan() {
			t := stdErrScanner.Text()
			if len(t) == 0 {
				continue
			}
			errb.WriteString(t)
			errb.WriteByte('\n')
			fmt.Fprintln(os.Stderr, t)
		}
	}()

	if err = cmd.Start(); err != nil {
		return nil, err
	}

	return &RunningProcess{
		cmd:          cmd,
		stdout:       outb,
		stderr:       errb,
		stdOutReader: stdOutReader,
		stdErrReader: stdErrReader,
	}, nil
}
