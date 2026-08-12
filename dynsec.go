package mosquitto

import (
	"context"
	"fmt"
)

type DynSec struct {
	m *Manager
}

func (m *Manager) DynamicSecurity() *DynSec {
	return &DynSec{m: m}
}

func (d *DynSec) Run(ctx context.Context, args ...string) (CommandResult, error) {
	if len(args) == 0 {
		return CommandResult{}, fmt.Errorf("dynamic security command is required")
	}
	c, cancel := d.m.withTimeout(ctx)
	defer cancel()
	base := []string{}
	if d.m.cfg.Host != "" {
		base = append(base, "-h", d.m.cfg.Host)
	}
	if d.m.cfg.Port != 0 {
		base = append(base, "-p", fmt.Sprint(d.m.cfg.Port))
	}
	base = append(base, args...)
	return d.m.ex.Run(c, d.m.cfg.CtrlPath, base...)
}
