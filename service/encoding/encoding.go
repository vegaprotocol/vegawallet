package encoding

import (
	"time"

	"go.uber.org/zap/zapcore"
)

// Duration is a wrapper over an actual duration so we can represent
// them as string in the toml configuration.
type Duration struct {
	time.Duration
}

// Get returns the stored duration.
func (d *Duration) Get() time.Duration {
	return d.Duration
}

// UnmarshalText unmarshal a duration from bytes.
func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

// MarshalText marshal a duration into bytes.
func (d Duration) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// LogLevel is wrapper over the actual log level
// so they can be specified as strings in the toml configuration.
type LogLevel struct {
	zapcore.Level
}

// Get return the store value.
func (l *LogLevel) Get() zapcore.Level {
	return l.Level
}

// UnmarshalText unmarshal a loglevel from bytes.
func (l *LogLevel) UnmarshalText(text []byte) error {
	return l.Level.UnmarshalText(text)
}

// MarshalText marshal a loglevel into bytes.
func (l LogLevel) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}
