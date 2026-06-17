package logger

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"

	"github.com/rs/zerolog"
)

type ZerologHandler struct {
	logger zerolog.Logger
	attrs  []slog.Attr
	groups []string
}

func NewZerologHandler(logger zerolog.Logger) *ZerologHandler {
	return &ZerologHandler{
		logger: logger,
		attrs:  make([]slog.Attr, 0),
		groups: make([]string, 0),
	}
}

func (h *ZerologHandler) Enabled(ctx context.Context, level slog.Level) bool {
	var zerologLevel zerolog.Level
	switch level {
	case slog.LevelDebug:
		zerologLevel = zerolog.DebugLevel
	case slog.LevelInfo:
		zerologLevel = zerolog.InfoLevel
	case slog.LevelWarn:
		zerologLevel = zerolog.WarnLevel
	case slog.LevelError:
		zerologLevel = zerolog.ErrorLevel
	default:
		zerologLevel = zerolog.InfoLevel
	}
	return h.logger.GetLevel() <= zerologLevel
}

func (h *ZerologHandler) Handle(ctx context.Context, record slog.Record) error {
	var event *zerolog.Event

	switch record.Level {
	case slog.LevelDebug:
		event = h.logger.Debug()
	case slog.LevelInfo:
		event = h.logger.Info()
	case slog.LevelWarn:
		event = h.logger.Warn()
	case slog.LevelError:
		event = h.logger.Error()
	default:
		event = h.logger.Info()
	}

	// Add attributes from the handler
	for _, attr := range h.attrs {
		h.addAttr(event, attr, h.groups)
	}

	// Add attributes from the record
	record.Attrs(func(attr slog.Attr) bool {
		h.addAttr(event, attr, h.groups)
		return true
	})

	if record.PC != 0 {
		frames := runtime.CallersFrames([]uintptr{record.PC})
		frame, _ := frames.Next()
		event.Str("caller", fmt.Sprintf("%s:%d", frame.File, frame.Line))
	}

	event.Msg(record.Message)
	return nil
}

func (h *ZerologHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandler := &ZerologHandler{
		logger: h.logger,
		attrs:  make([]slog.Attr, len(h.attrs)+len(attrs)),
		groups: make([]string, len(h.groups)),
	}
	copy(newHandler.attrs, h.attrs)
	copy(newHandler.attrs[len(h.attrs):], attrs)
	copy(newHandler.groups, h.groups)
	return newHandler
}

func (h *ZerologHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	newHandler := &ZerologHandler{
		logger: h.logger,
		attrs:  make([]slog.Attr, len(h.attrs)),
		groups: make([]string, len(h.groups)+1),
	}
	copy(newHandler.attrs, h.attrs)
	copy(newHandler.groups, h.groups)
	newHandler.groups[len(h.groups)] = name
	return newHandler
}

func (h *ZerologHandler) addAttr(event *zerolog.Event, attr slog.Attr, groups []string) {
	key := attr.Key
	if len(groups) > 0 {
		prefix := ""
		for _, g := range groups {
			prefix += g + "."
		}
		key = prefix + key
	}

	value := attr.Value.Any()

	switch v := value.(type) {
	case string:
		event.Str(key, v)
	case int:
		event.Int(key, v)
	case int64:
		event.Int64(key, v)
	case float64:
		event.Float64(key, v)
	case bool:
		event.Bool(key, v)
	case error:
		event.Err(v)
	default:
		event.Interface(key, v)
	}
}
