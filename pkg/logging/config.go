package logging

import (
	"github.com/dexlabsio/garlic/pkg/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Level    string `json:"level" mapstructure:"level" yaml:"level"`
	Encoding string `json:"encoding" mapstructure:"encoding" yaml:"encoding"`
}

func (c *Config) Parse() *zap.Config {
	lvl, err := zap.ParseAtomicLevel(c.Level)
	if err != nil {
		panic(err)
	}

	zapConfig := &zap.Config{
		// Level is the minimum enabled logging level. Note that this is a dynamic
		// level, so calling Config.Level.SetLevel will atomically change the log
		// level of all loggers descended from this config.
		Level: lvl,
		// Development puts the logger in development mode, which changes the
		// behavior of DPanicLevel and takes stacktraces more liberally.
		Development: true,
		// DisableCaller stops annotating logs with the calling function's file
		// name and line number. By default, all logs are annotated.
		DisableCaller: false,
		// DisableStacktrace completely disables automatic stacktrace capturing. By
		// default, stacktraces are captured for WarnLevel and above logs in
		// development and ErrorLevel and above in production.
		DisableStacktrace: false,
		// Sampling sets a sampling policy. A nil SamplingConfig disables sampling.
		Sampling: nil,
		// Encoding sets the logger's encoding. Valid values are "json" and
		// "console", as well as any third-party encodings registered via
		// RegisterEncoder.
		Encoding: c.Encoding,
		// EncoderConfig sets options for the chosen encoder. See
		// zapcore.EncoderConfig for details.
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "msg",
			LevelKey:       "level",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		},

		// OutputPaths is a list of URLs or file paths to write logging output to.
		// See Open for details.
		OutputPaths: []string{"stdout"},
		// ErrorOutputPaths is a list of URLs to write internal logger errors to.
		// The default is standard error.
		//
		// Note that this setting only affects internal errors; for sample code that
		// sends error-level logs to a different location from info- and debug-level
		// logs, see the package-level AdvancedConfiguration example.
		ErrorOutputPaths: []string{"stderr"},
		// InitialFields is a collection of fields to add to the root logger.
		InitialFields: map[string]interface{}{
			"version": global.Version,
		},
	}

	// In development mode, we want to have a more human-readable output
	if c.Encoding == "console" {
		zapConfig.EncoderConfig = zapcore.EncoderConfig{
			MessageKey: "msg",
		}

		zapConfig.InitialFields = map[string]interface{}{}
	}

	return zapConfig
}

func Defaults() *Config {
	return &Config{
		Level:    "error",
		Encoding: "json",
	}
}
