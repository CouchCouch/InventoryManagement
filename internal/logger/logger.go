package logger

import (
	"log/slog"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

func Initialize(logLevel string, isDevelopment bool) *slog.Logger {
	// Parse log level
	var zerologLevel zerolog.Level
	switch strings.ToLower(logLevel) {
	case "debug":
		zerologLevel = zerolog.DebugLevel
	case "info":
		zerologLevel = zerolog.InfoLevel
	case "warn", "warning":
		zerologLevel = zerolog.WarnLevel
	case "error":
		zerologLevel = zerolog.ErrorLevel
	default:
		zerologLevel = zerolog.InfoLevel
	}

	// Configure zerolog output format based on environment
	var zerologLogger zerolog.Logger
	if isDevelopment {
		// Console output with human-readable formatting
		zerologLogger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}).
			Level(zerologLevel).
			With().
			Timestamp().
			Logger()
	} else {
		// JSON output for production
		zerologLogger = zerolog.New(os.Stdout).
			Level(zerologLevel).
			With().
			Timestamp().
			Logger()
	}

	// Create slog logger with zerolog handler
	handler := NewZerologHandler(zerologLogger)
	return slog.New(handler)
}
