package cmd

import (
	"fmt"
	"os"
	"time"

	"code.vegaprotocol.io/shared/paths"
	"code.vegaprotocol.io/vegawallet/cmd/flags"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var SupportedLogLevels = []string{
	zapcore.DebugLevel.String(),
	zapcore.InfoLevel.String(),
	zapcore.WarnLevel.String(),
	zapcore.ErrorLevel.String(),
}

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

func BuildJSONLogger(level string, vegaPaths paths.Paths, logsDir paths.StatePath) (*zap.Logger, string, error) {
	cfg := DefaultConfig()

	l, err := getLevel(level)
	if err != nil {
		return nil, "", err
	}

	cfg.Level = zap.NewAtomicLevelAt(*l)

	pid := os.Getpid()
	date := time.Now().UTC().Format("2006-01-02-15-04-05")
	pathSuffix := fmt.Sprintf("%d-%s.log", pid, date)

	logFile := paths.JoinStatePath(logsDir, pathSuffix)
	appLogPath, err := vegaPaths.CreateStatePathFor(logFile)
	if err != nil {
		return nil, "", fmt.Errorf("failed getting path for %s: %w", logFile, err)
	}

	zapLogPath := toZapLogPath(appLogPath)
	cfg.OutputPaths = []string{zapLogPath}
	cfg.ErrorOutputPaths = []string{zapLogPath}

	log, err := cfg.Build()
	if err != nil {
		return nil, "", fmt.Errorf("couldn't create logger: %w", err)
	}
	return log, appLogPath, nil
}

func BuildLogger(output, level string) (*zap.Logger, error) {
	cfg := DefaultConfig()

	l, err := getLevel(level)
	if err != nil {
		return nil, err
	}

	cfg.Level = zap.NewAtomicLevelAt(*l)

	if output == flags.InteractiveOutput {
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cfg.Encoding = "console"
	}

	log, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("couldn't create logger: %w", err)
	}
	return log, nil
}

func ValidateLogLevel(level string) error {
	if isSupportedLogLevel(level) {
		return nil
	}

	return NewUnsupportedFlagValueError(level)
}

func NewUnsupportedFlagValueError(level string) error {
	supportedLogLevels := make([]interface{}, len(SupportedLogLevels))
	for i := range SupportedLogLevels {
		supportedLogLevels[i] = SupportedLogLevels[i]
	}

	return flags.UnsupportedFlagValueError("level", level, supportedLogLevels)
}

func getLevel(level string) (*zapcore.Level, error) {
	if !isSupportedLogLevel(level) {
		return nil, UnsupportedLoggerLevelError(level)
	}

	l := new(zapcore.Level)

	err := l.UnmarshalText([]byte(level))
	if err != nil {
		return nil, fmt.Errorf("couldn't parse logger level: %w", err)
	}

	return l, nil
}

func isSupportedLogLevel(level string) bool {
	for _, supported := range SupportedLogLevels {
		if level == supported {
			return true
		}
	}
	return false
}
