package logging

import "go.uber.org/zap"

var Global *zap.Logger

// Init initializes a new global logger with a standard configuration
// prepared for logging stacks like ELK. Users can still define the
// log level using strings. Avalable options: "warn", "debug", "info" and "error".
func Init(config *Config) {
	if Global != nil {
		Global.Fatal("Failed to configure a new Global: this is already set")
	}

	Global = zap.Must(config.Parse().Build())
}
