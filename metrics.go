package mosquitto

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

type Metrics struct {
	ConnectedClients int
	TotalClients     int
	MessagesReceived uint64
	MessagesSent     uint64
	BytesReceived    uint64
	BytesSent        uint64
}

func (m *Manager) ConnectedClientCount(ctx context.Context) (int, error) {
	values, err := m.sysMetrics(ctx, []string{"$SYS/broker/clients/connected"})
	if err != nil {
		return 0, err
	}
	return parseIntMetric(values["$SYS/broker/clients/connected"])
}

func (m *Manager) Metrics(ctx context.Context) (Metrics, error) {
	values, err := m.sysMetrics(ctx, []string{
		"$SYS/broker/clients/connected",
		"$SYS/broker/clients/total",
		"$SYS/broker/messages/received",
		"$SYS/broker/messages/sent",
		"$SYS/broker/bytes/received",
		"$SYS/broker/bytes/sent",
	})
	if err != nil {
		return Metrics{}, err
	}
	return Metrics{
		ConnectedClients: mustInt(values["$SYS/broker/clients/connected"]),
		TotalClients:     mustInt(values["$SYS/broker/clients/total"]),
		MessagesReceived: mustUint(values["$SYS/broker/messages/received"]),
		MessagesSent:     mustUint(values["$SYS/broker/messages/sent"]),
		BytesReceived:    mustUint(values["$SYS/broker/bytes/received"]),
		BytesSent:        mustUint(values["$SYS/broker/bytes/sent"]),
	}, nil
}

func (m *Manager) sysMetrics(ctx context.Context, topics []string) (map[string]string, error) {
	c, cancel := m.withTimeout(ctx)
	defer cancel()
	args := []string{"-h", m.cfg.Host, "-p", fmt.Sprint(m.cfg.Port), "-C", fmt.Sprint(len(topics)), "-v"}
	for _, t := range topics {
		args = append(args, "-t", t)
	}
	r, err := m.ex.Run(c, m.cfg.SubPath, args...)
	if err != nil {
		return nil, err
	}
	values := make(map[string]string)
	for _, line := range strings.Split(r.Stdout, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 2 {
			values[parts[0]] = parts[1]
		}
	}
	return values, nil
}

func parseIntMetric(v string) (int, error) {
	n, err := strconv.Atoi(strings.TrimSpace(v))
	if err != nil {
		return 0, fmt.Errorf("invalid metric %q: %w", v, err)
	}
	return n, nil
}
func mustInt(v string) int {
	n, _ := parseIntMetric(v)
	return n
}
func mustUint(v string) uint64 {
	n, _ := strconv.ParseUint(strings.TrimSpace(v), 10, 64)
	return n
}
