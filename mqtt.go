package mosquitto

import (
	"context"
	"fmt"
	"strings"
)

type Message struct {
	Topic    string
	Payload  []byte
	QoS      int
	Retained bool
}

type PublishOptions struct {
	QoS      int
	Retain   bool
	Username string
	Password string
	ClientID string
	Protocol string
	Timeout  int
}

func (m *Manager) Publish(ctx context.Context, msg Message, opts PublishOptions) error {
	if msg.Topic == "" {
		return fmt.Errorf("topic is required")
	}
	if opts.QoS < 0 || opts.QoS > 2 {
		return fmt.Errorf("qos must be 0, 1, or 2")
	}
	c, cancel := m.withTimeout(ctx)
	defer cancel()

	args := []string{"-h", m.cfg.Host, "-p", fmt.Sprint(m.cfg.Port), "-t", msg.Topic, "-m", string(msg.Payload), "-q", fmt.Sprint(opts.QoS)}
	if msg.Retained || opts.Retain {
		args = append(args, "-r")
	}
	if opts.Username != "" {
		args = append(args, "-u", opts.Username)
	}
	if opts.Password != "" {
		args = append(args, "-P", opts.Password)
	}
	if opts.ClientID != "" {
		args = append(args, "-i", opts.ClientID)
	}
	if opts.Protocol != "" {
		args = append(args, "-V", opts.Protocol)
	}
	_, err := m.ex.Run(c, m.cfg.PubPath, args...)
	return err
}

type SubscribeOptions struct {
	QoS      int
	Username string
	Password string
	ClientID string
	Protocol string
	Verbose  bool
}

type Subscription struct {
	Command *Process
}

func (m *Manager) Subscribe(ctx context.Context, topic string, opts SubscribeOptions) (*Subscription, error) {
	if topic == "" {
		return nil, fmt.Errorf("topic is required")
	}
	if opts.QoS < 0 || opts.QoS > 2 {
		return nil, fmt.Errorf("qos must be 0, 1, or 2")
	}

	args := []string{"-h", m.cfg.Host, "-p", fmt.Sprint(m.cfg.Port), "-t", topic, "-q", fmt.Sprint(opts.QoS)}
	if opts.Username != "" {
		args = append(args, "-u", opts.Username)
	}
	if opts.Password != "" {
		args = append(args, "-P", opts.Password)
	}
	if opts.ClientID != "" {
		args = append(args, "-i", opts.ClientID)
	}
	if opts.Protocol != "" {
		args = append(args, "-V", opts.Protocol)
	}
	if opts.Verbose {
		args = append(args, "-v")
	}

	p, err := NewProcess(ctx, m.cfg.SubPath, args...)
	if err != nil {
		return nil, err
	}
	return &Subscription{Command: p}, nil
}

func (m *Manager) Request(ctx context.Context, requestTopic, payload, responseTopic string, opts PublishOptions) ([]byte, error) {
	if requestTopic == "" || responseTopic == "" {
		return nil, fmt.Errorf("request and response topics are required")
	}
	c, cancel := m.withTimeout(ctx)
	defer cancel()

	args := []string{
		"-h", m.cfg.Host, "-p", fmt.Sprint(m.cfg.Port),
		"-t", requestTopic, "-m", payload,
		"-e", responseTopic,
	}
	if opts.QoS >= 0 && opts.QoS <= 2 {
		args = append(args, "-q", fmt.Sprint(opts.QoS))
	}
	if opts.Username != "" {
		args = append(args, "-u", opts.Username)
	}
	if opts.Password != "" {
		args = append(args, "-P", opts.Password)
	}
	r, err := m.ex.Run(c, m.cfg.RRPath, args...)
	if err != nil {
		return nil, err
	}
	return []byte(strings.TrimSuffix(r.Stdout, "\r\n")), nil
}
