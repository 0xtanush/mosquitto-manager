package mosquitto

import (
	"context"
	"errors"
	"os/user"
	"runtime"
	"time"
)

var (
	ErrUnsupported = errors.New("operation is not supported on this platform")
	ErrNotRunning  = errors.New("mosquitto broker is not running")
)

type Config struct {
	MosquittoPath string
	PubPath       string
	SubPath       string
	RRPath        string
	CtrlPath      string
	PasswdPath    string
	SignalPath    string

	ServiceName  string
	ConfigFile   string
	PasswordFile string

	Host     string
	Port     int
	Username string
	Password string

	CommandTimeout time.Duration
	ServiceTimeout time.Duration
}

func DefaultConfig() Config {
	c := Config{
		MosquittoPath:  "mosquitto",
		PubPath:        "mosquitto_pub",
		SubPath:        "mosquitto_sub",
		RRPath:         "mosquitto_rr",
		CtrlPath:       "mosquitto_ctrl",
		PasswdPath:     "mosquitto_passwd",
		SignalPath:     "mosquitto_signal",
		ServiceName:    "mosquitto",
		ConfigFile:     "",
		PasswordFile:   "",
		Host:           "localhost",
		Port:           1883,
		CommandTimeout: 15 * time.Second,
		ServiceTimeout: 30 * time.Second,
	}
	if runtime.GOOS == "windows" {
		// Keep executable names extension-free; exec.LookPath resolves .exe.
	}
	return c
}

type Manager struct {
	cfg Config
	ex  Executor
}

func New(cfg Config) *Manager {
	if cfg.MosquittoPath == "" {
		cfg.MosquittoPath = "mosquitto"
	}
	if cfg.PubPath == "" {
		cfg.PubPath = "mosquitto_pub"
	}
	if cfg.SubPath == "" {
		cfg.SubPath = "mosquitto_sub"
	}
	if cfg.RRPath == "" {
		cfg.RRPath = "mosquitto_rr"
	}
	if cfg.CtrlPath == "" {
		cfg.CtrlPath = "mosquitto_ctrl"
	}
	if cfg.PasswdPath == "" {
		cfg.PasswdPath = "mosquitto_passwd"
	}
	if cfg.SignalPath == "" {
		cfg.SignalPath = "mosquitto_signal"
	}
	if cfg.ServiceName == "" {
		cfg.ServiceName = "mosquitto"
	}
	if cfg.Host == "" {
		cfg.Host = "localhost"
	}
	if cfg.Port == 0 {
		cfg.Port = 1883
	}
	if cfg.CommandTimeout <= 0 {
		cfg.CommandTimeout = 15 * time.Second
	}
	if cfg.ServiceTimeout <= 0 {
		cfg.ServiceTimeout = 30 * time.Second
	}
	return &Manager{cfg: cfg, ex: NewOSExecutor()}
}

func (m *Manager) Config() Config { return m.cfg }

func (m *Manager) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithTimeout(ctx, m.cfg.CommandTimeout)
}

func (m *Manager) withServiceTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithTimeout(ctx, m.cfg.ServiceTimeout)
}

func currentOS() string { return runtime.GOOS }
func currentUser() string {
	if u, err := user.Current(); err == nil {
		return u.Username
	}
	return ""
}
