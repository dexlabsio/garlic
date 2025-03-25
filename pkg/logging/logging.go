package logging

import "go.uber.org/zap"

// New initializes a new global logger with a standard configuration
// prepared for logging stacks like ELK. Users can still define the
// log level using strings. Avalable options: "warn", "debug", "info" and "error".
func New(config *Config) *zap.Logger {
	return zap.Must(config.Parse().Build())
}
