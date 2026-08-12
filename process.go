package mosquitto

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"sync"
)

type Process struct {
	Cmd    *exec.Cmd
	Stdout io.ReadCloser
	Stderr io.ReadCloser

	mu sync.Mutex
}

func NewProcess(ctx context.Context, command string, args ...string) (*Process, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	cmd := exec.CommandContext(ctx, command, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		_ = stdout.Close()
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		_ = stdout.Close()
		_ = stderr.Close()
		return nil, err
	}
	return &Process{Cmd: cmd, Stdout: stdout, Stderr: stderr}, nil
}

func (p *Process) Wait() error { return p.Cmd.Wait() }

func (p *Process) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.Cmd.Process == nil {
		return nil
	}
	return p.Cmd.Process.Kill()
}

func (p *Process) Lines() <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		sc := bufio.NewScanner(p.Stdout)
		for sc.Scan() {
			ch <- sc.Text()
		}
	}()
	return ch
}

func (p *Process) ReadAll() ([]byte, error) {
	var b bytes.Buffer
	if p.Stdout != nil {
		_, _ = io.Copy(&b, p.Stdout)
	}
	return b.Bytes(), nil
}

func (p *Process) PID() int {
	if p.Cmd == nil || p.Cmd.Process == nil {
		return 0
	}
	return p.Cmd.Process.Pid
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
