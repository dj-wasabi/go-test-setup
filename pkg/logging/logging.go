package logging

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"

	"werner-dijkerman.nl/test-setup/pkg/config"
)

// Logger interface defines methods for logging at different levels.
// type LoggerInterface interface {
// 	Debug(string, string, ...any)
// 	Info(string, string, ...any)
// 	Warn(string, string, ...any)
// 	Error(string, string, ...any)
// }

// type CustomLogger struct {
// 	LoggerInterface
// 	logger *slog.Logger
// }

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

// // Debug logs a debug message.
// func (l *CustomLogger) Debug(logid, msg string, args ...any) {
// 	l.logger.Debug("log_id", logid, fmt.Sprintf(msg, args...))
// }

// // Info logs an info message.
// func (l *CustomLogger) Info(logid, msg string, args ...any) {
// 	l.logger.Info("log_id", logid, fmt.Sprintf(msg, args...))
// }

// // Warn logs a warning message.
// func (l *CustomLogger) Warn(logid, msg string, args ...any) {
// 	l.logger.Warn("log_id", logid, fmt.Sprintf(msg, args...))
// }

// // Error logs an error message.
// func (l *CustomLogger) Error(logid, msg string, args ...any) {
// 	l.logger.Error("log_id", logid, fmt.Sprintf(msg, args...))
// }
