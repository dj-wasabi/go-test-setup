package logging

import (
	"log/slog"
	"testing"
)

func Test_parse_loglevel_upper(t *testing.T) {
	loglevel := "DEBUG"
	parsed_log := parseLogLevel(loglevel)
	if parsed_log != slog.LevelDebug {
		t.Errorf("Expected log level to be Debug, got %v", parsed_log)
	}
}

func Test_parse_loglevel_lower(t *testing.T) {
	loglevel := "warning"
	parsed_log := parseLogLevel(loglevel)
	if parsed_log != slog.LevelWarn {
		t.Errorf("Expected log level to be Warning, got %v", parsed_log)
	}
}
