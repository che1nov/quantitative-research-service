package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

// Logger описывает контракт логгера для всех слоев приложения.
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	DebugContext(ctx context.Context, msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
	With(args ...any) Logger
}

type slogLogger struct {
	logger *slog.Logger
}

// New создает единый логгер приложения.
func New(level string) Logger {
	var logLevel slog.Level
	switch strings.ToLower(level) {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	return &slogLogger{logger: slog.New(handler)}
}

func (l *slogLogger) Debug(msg string, args ...any) { l.logger.Debug(msg, args...) }
func (l *slogLogger) Info(msg string, args ...any)  { l.logger.Info(msg, args...) }
func (l *slogLogger) Warn(msg string, args ...any)  { l.logger.Warn(msg, args...) }
func (l *slogLogger) Error(msg string, args ...any) { l.logger.Error(msg, args...) }

func (l *slogLogger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.logger.DebugContext(ctx, msg, args...)
}

func (l *slogLogger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, args...)
}

func (l *slogLogger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.logger.WarnContext(ctx, msg, args...)
}

func (l *slogLogger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, args...)
}

func (l *slogLogger) With(args ...any) Logger {
	return &slogLogger{logger: l.logger.With(args...)}
}
