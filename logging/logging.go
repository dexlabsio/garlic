package logging

import "go.uber.org/zap"

var singleton *zap.Logger

// Init initializes a new global logger with a standard configuration
// prepared for logging stacks like ELK. Users can still define the
// log level using strings. Avalable options: "warn", "debug", "info" and "error".
func Init(config *Config) {
	if singleton != nil {
		singleton.Fatal("Failed to configure a new Global: this is already set")
	}

	singleton = zap.Must(config.Parse().Build())
}

func Global() *zap.Logger {
	if singleton == nil {
		singleton = zap.Must(Defaults().Parse().Build())
	}

	return singleton
}
