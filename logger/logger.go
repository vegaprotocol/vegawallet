package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var SupportedLogLevels = []zapcore.Level{zapcore.DebugLevel, zapcore.InfoLevel, zapcore.WarnLevel, zapcore.ErrorLevel}

type LoggerError struct {
	message string
}

func UnsupportedLoggerLevelError(l string) LoggerError {
	return LoggerError{
		message: fmt.Sprintf("unsupported logger level \"%s\"", l),
	}
}

func (e LoggerError) Error() string {
	return e.message
}

func DefaultConfig() zap.Config {
	return zap.Config{
		Level:    zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "message",
			LevelKey:       "level",
			TimeKey:        "@timestamp",
			NameKey:        "logger",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			LineEnding:     "\n",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		},
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		DisableStacktrace: true,
	}
}

func Build(output, level string) (*zap.Logger, error) {
	cfg := DefaultConfig()

	l, err := getLevel(level)
	if err != nil {
		return nil, err
	}

	cfg.Level = zap.NewAtomicLevelAt(*l)

	if output == "human" {
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cfg.Encoding = "console"
	}

	log, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("couldn't create logger: %w", err)
	}
	return log, nil
}

func getLevel(level string) (*zapcore.Level, error) {
	l := new(zapcore.Level)

	err := l.UnmarshalText([]byte(level))
	if err != nil {
		return nil, fmt.Errorf("couldn't parse logger level: %w", err)
	}

	if !isSupportedLogLevel(*l) {
		return nil, UnsupportedLoggerLevelError(l.String())
	}

	return l, nil
}

func isSupportedLogLevel(selected zapcore.Level) bool {
	for _, supported := range SupportedLogLevels {
		if selected == supported {
			return true
		}
	}
	return false
}
