package options

import (
	"io"
	"log/slog"
)

type Option func(*Config)

func WithLogLevel(level slog.Level) Option {
	return func(cfg *Config) {
		cfg.log.level = level
	}
}

func WithLogWriter(writer io.Writer) Option {
	return func(cfg *Config) {
		cfg.log.writer = writer
	}
}
