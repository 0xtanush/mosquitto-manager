# mosquitto-manager v0.2

A small cross-platform Go control library for an installed Eclipse Mosquitto broker.

## Platforms

- Linux/POSIX: service control through `systemctl` when configured, signals through `mosquitto_signal` (preferred) or `kill` for reload.
- Windows: Windows Service Control Manager through `sc.exe`, and `mosquitto_signal.exe` for broker signals.
- MQTT CLI operations: `mosquitto_pub`, `mosquitto_sub`, and `mosquitto_rr`.
- Security: `mosquitto_passwd` and `mosquitto_ctrl`.

The library intentionally wraps Mosquitto's existing tools rather than reimplementing the broker.

## Important

This package assumes Mosquitto is already installed. It does not install Mosquitto or register a Windows service for you.

Windows service operations require the Mosquitto service name configured in `Config.ServiceName` (default: `mosquitto`).

Run the tests with:

```text
go test ./...
```

The integration tests are intentionally not included because they require a real Mosquitto installation and service configuration.
"# mosquitto-manager" 
