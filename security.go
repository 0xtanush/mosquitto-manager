package mosquitto

import (
	"context"
	"fmt"
	"os"
	"strings"
)

type PasswordUser struct {
	Username string
}

func (m *Manager) CreateUser(ctx context.Context, username, password string) error {
	if username == "" {
		return fmt.Errorf("username is required")
	}
	if m.cfg.PasswordFile == "" {
		return fmt.Errorf("password file is not configured")
	}
	c, cancel := m.withTimeout(ctx)
	defer cancel()
	_, err := m.ex.Run(c, m.cfg.PasswdPath, m.cfg.PasswordFile, username, password)
	return err
}

func (m *Manager) CreatePasswordFileUser(ctx context.Context, username, password string) error {
	if m.cfg.PasswordFile == "" {
		return fmt.Errorf("password file is not configured")
	}
	if _, err := os.Stat(m.cfg.PasswordFile); os.IsNotExist(err) {
		c, cancel := m.withTimeout(ctx)
		defer cancel()
		_, err = m.ex.Run(c, m.cfg.PasswdPath, "-c", m.cfg.PasswordFile, username, password)
		return err
	}
	return m.CreateUser(ctx, username, password)
}

func (m *Manager) DeleteUser(ctx context.Context, username string) error {
	if m.cfg.PasswordFile == "" {
		return fmt.Errorf("password file is not configured")
	}
	c, cancel := m.withTimeout(ctx)
	defer cancel()
	_, err := m.ex.Run(c, m.cfg.PasswdPath, "-D", m.cfg.PasswordFile, username)
	return err
}

func (m *Manager) ListPasswordFileUsers(ctx context.Context) ([]PasswordUser, error) {
	if m.cfg.PasswordFile == "" {
		return nil, fmt.Errorf("password file is not configured")
	}
	data, err := os.ReadFile(m.cfg.PasswordFile)
	if err != nil {
		return nil, err
	}
	var users []PasswordUser
	for _, line := range strings.Split(string(data), "\n") {
		if i := strings.IndexByte(line, ':'); i > 0 {
			users = append(users, PasswordUser{Username: line[:i]})
		}
	}
	return users, nil
}
