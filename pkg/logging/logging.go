package logging

import (
	"log/slog"
	"os"
	"strings"
	"sync"

	"werner-dijkerman.nl/test-setup/pkg/config"
)

var (
	once   sync.Once
	logger *slog.Logger
)

func parseLogLevel(logLevelString string) slog.Level {
	switch strings.ToUpper(logLevelString) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "WARNING":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func Load() *slog.Logger {
	c := config.ReadConfig()
	logLevel := parseLogLevel(c.Logging.Level)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))

	return logger
}

func Initialize() *slog.Logger {
	once.Do(func() {
		logger = Load()
	})
	return logger
}
