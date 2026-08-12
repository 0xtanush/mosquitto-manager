package mosquitto

import (
	"context"
	"fmt"
	"os"
)

func (m *Manager) ValidateConfig(ctx context.Context) error {
	if m.cfg.ConfigFile == "" {
		return fmt.Errorf("config file is not configured")
	}
	c, cancel := m.withTimeout(ctx)
	defer cancel()
	_, err := m.ex.Run(c, m.cfg.MosquittoPath, "-c", m.cfg.ConfigFile, "-t")
	return err
}

func (m *Manager) ReadConfig() ([]byte, error) {
	if m.cfg.ConfigFile == "" {
		return nil, fmt.Errorf("config file is not configured")
	}
	return os.ReadFile(m.cfg.ConfigFile)
}
