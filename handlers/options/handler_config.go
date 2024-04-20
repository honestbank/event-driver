package options

import (
	"io"
	"log/slog"
	"os"
)

type Config struct {
	log LogConfig
}

type LogConfig struct {
	level  slog.Level
	writer io.Writer
}

func DefaultOptions() Config {
	return Config{
		log: LogConfig{
			level:  slog.LevelInfo,
			writer: os.Stdout,
		},
	}
}

func (c Config) GetLogLevel() slog.Level {
	return c.log.level
}

func (c Config) GetLogWriter() io.Writer {
	return c.log.writer
}
