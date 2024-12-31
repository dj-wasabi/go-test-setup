package logging

import (
	"fmt"
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

// parseLogLevel parses the log level from a string.
func parseLogLevel(logLevelString string) slog.Level {
	switch strings.ToUpper(logLevelString) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN", "WARNING":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// Load initializes and returns a new slog.Logger.
func Load() *slog.Logger {
	c := config.ReadConfig()
	logLevel := parseLogLevel(c.Logging.Level)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))

	logger.Info(fmt.Sprintf("The loglevel '%v' is active.", logLevel))
	return logger
}

// Initialize initializes the custom logger singleton.
func Initialize() *slog.Logger {
	once.Do(func() {
		logger = Load()
	})
	return logger
}
