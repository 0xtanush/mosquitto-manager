//go:build windows

package mosquitto

// Windows-specific behavior is intentionally implemented through the native
// Windows Service Control Manager (sc.exe) and Mosquitto's mosquitto_signal.
// Keeping these operations behind Manager means callers do not need to know
// Windows command details.
