package mosquitto

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

type BrokerStatus struct {
	Running   bool
	Service   string
	PID       int
	RawOutput string
}

func (m *Manager) Start(ctx context.Context) error {
	switch currentOS() {
	case "windows":
		return m.windowsService(ctx, "start")
	default:
		return m.posixService(ctx, "start")
	}
}

func (m *Manager) Stop(ctx context.Context) error {
	switch currentOS() {
	case "windows":
		return m.windowsService(ctx, "stop")
	default:
		return m.posixService(ctx, "stop")
	}
}

func (m *Manager) Restart(ctx context.Context) error {
	switch currentOS() {
	case "windows":
		if err := m.windowsService(ctx, "stop"); err != nil {
			// If already stopped, continue to start. A later status check
			// will determine whether the broker actually came up.
		}
		return m.windowsService(ctx, "start")
	default:
		return m.posixService(ctx, "restart")
	}
}

func (m *Manager) Reload(ctx context.Context) error {
	// mosquitto_signal is the portable Mosquitto-native mechanism.
	// On Windows it is the documented way to signal the broker.
	// On POSIX, SIGHUP is also possible, but using mosquitto_signal keeps
	// the library's behavior aligned across platforms.
	if m.cfg.SignalPath == "" {
		return ErrUnsupported
	}
	c, cancel := m.withTimeout(ctx)
	defer cancel()

	_, err := m.ex.Run(c, m.cfg.SignalPath, "-a", "config-reload")
	return err
}

func (m *Manager) Status(ctx context.Context) (BrokerStatus, error) {
	if currentOS() == "windows" {
		return m.windowsStatus(ctx)
	}
	return m.posixStatus(ctx)
}

func (m *Manager) IsRunning(ctx context.Context) (bool, error) {
	s, err := m.Status(ctx)
	return s.Running, err
}

func (m *Manager) Version(ctx context.Context) (string, error) {
	c, cancel := m.withTimeout(ctx)
	defer cancel()
	r, err := m.ex.Run(c, m.cfg.MosquittoPath, "-h")
	if err != nil {
		return "", err
	}
	text := r.Stdout + "\n" + r.Stderr
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if strings.Contains(strings.ToLower(line), "mosquitto version") {
			return line, nil
		}
	}
	return strings.TrimSpace(text), nil
}

func (m *Manager) posixService(ctx context.Context, action string) error {
	c, cancel := m.withServiceTimeout(ctx)
	defer cancel()
	_, err := m.ex.Run(c, "systemctl", action, m.cfg.ServiceName)
	return err
}

func (m *Manager) windowsService(ctx context.Context, action string) error {
	c, cancel := m.withServiceTimeout(ctx)
	defer cancel()

	// sc.exe is part of Windows and is the stable command-line interface
	// to the Windows Service Control Manager.
	_, err := m.ex.Run(c, "sc.exe", action, m.cfg.ServiceName)
	return err
}

func (m *Manager) posixStatus(ctx context.Context) (BrokerStatus, error) {
	c, cancel := m.withServiceTimeout(ctx)
	defer cancel()
	r, err := m.ex.Run(c, "systemctl", "show", m.cfg.ServiceName,
		"--property=ActiveState,MainPID", "--no-page")
	if err != nil {
		return BrokerStatus{Service: m.cfg.ServiceName, RawOutput: r.Stdout + r.Stderr}, err
	}

	s := BrokerStatus{Service: m.cfg.ServiceName, RawOutput: r.Stdout}
	var active string
	for _, line := range strings.Split(r.Stdout, "\n") {
		k, v, ok := strings.Cut(strings.TrimSpace(line), "=")
		if !ok {
			continue
		}
		switch k {
		case "ActiveState":
			active = v
		case "MainPID":
			s.PID, _ = strconv.Atoi(v)
		}
	}
	s.Running = active == "active"
	return s, nil
}

func (m *Manager) windowsStatus(ctx context.Context) (BrokerStatus, error) {
	c, cancel := m.withServiceTimeout(ctx)
	defer cancel()
	r, err := m.ex.Run(c, "sc.exe", "query", m.cfg.ServiceName)
	raw := r.Stdout + r.Stderr
	s := BrokerStatus{Service: m.cfg.ServiceName, RawOutput: raw}
	if err != nil {
		return s, err
	}

	for _, line := range strings.Split(r.Stdout, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "STATE") {
			// Example: STATE : 4 RUNNING
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				s.Running = len(fields) >= 4 && strings.EqualFold(fields[3], "RUNNING")
			}
		}
	}
	return s, nil
}

func (m *Manager) EnsureRunning(ctx context.Context) error {
	s, err := m.Status(ctx)
	if err != nil {
		return err
	}
	if s.Running {
		return nil
	}
	if err := m.Start(ctx); err != nil {
		return err
	}
	s, err = m.Status(ctx)
	if err != nil {
		return err
	}
	if !s.Running {
		return fmt.Errorf("mosquitto service did not reach running state")
	}
	return nil
}
