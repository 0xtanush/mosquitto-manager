package mosquitto

import (
	"context"
	"strings"
	"testing"
)

type fakeExecutor struct {
	result CommandResult
	err    error
	calls  [][]string
}

func (f *fakeExecutor) Run(_ context.Context, command string, args ...string) (CommandResult, error) {
	f.calls = append(f.calls, append([]string{command}, args...))
	return f.result, f.err
}

func TestDefaultConfig(t *testing.T) {
	c := DefaultConfig()
	if c.ServiceName != "mosquitto" {
		t.Fatal(c.ServiceName)
	}
}

func TestCommandResultString(t *testing.T) {
	r := CommandResult{Command: "sc.exe", Args: []string{"query", "mosquitto"}, ExitCode: 0}
	if !strings.Contains(r.String(), "sc.exe") {
		t.Fatal(r.String())
	}
}
