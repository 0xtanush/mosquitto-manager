package mosquitto

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

type CommandResult struct {
	Command  string
	Args     []string
	Stdout   string
	Stderr   string
	ExitCode int
}

func (r CommandResult) String() string {
	return fmt.Sprintf("%s %s (exit=%d)", r.Command, strings.Join(r.Args, " "), r.ExitCode)
}

type Executor interface {
	Run(ctx context.Context, command string, args ...string) (CommandResult, error)
}

type OSExecutor struct{}

func NewOSExecutor() Executor { return OSExecutor{} }

func (OSExecutor) Run(ctx context.Context, command string, args ...string) (CommandResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	cmd := exec.CommandContext(ctx, command, args...)
	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	err := cmd.Run()
	result := CommandResult{
		Command:  command,
		Args:     append([]string(nil), args...),
		Stdout:   out.String(),
		Stderr:   errOut.String(),
		ExitCode: 0,
	}
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1
		}
		return result, fmt.Errorf("%s: %w", result.String(), err)
	}
	_ = runtime.GOOS
	return result, nil
}
