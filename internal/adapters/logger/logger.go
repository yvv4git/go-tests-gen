package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/yvv4git/go-tests-gen/internal/ports"
)

type SlogAdapter struct {
	logger *slog.Logger
}

func NewSlogAdapter(l *slog.Logger) *SlogAdapter {
	return &SlogAdapter{logger: l}
}

func (a *SlogAdapter) Debug(msg string, args ...any) {
	a.logger.Debug(msg, args...)
}

func (a *SlogAdapter) Info(msg string, args ...any) {
	a.logger.Info(msg, args...)
}

func (a *SlogAdapter) Warn(msg string, args ...any) {
	a.logger.Warn(msg, args...)
}

func (a *SlogAdapter) Error(msg string, args ...any) {
	a.logger.Error(msg, args...)
}

func (a *SlogAdapter) DebugContext(ctx context.Context, msg string, args ...any) {
	a.logger.DebugContext(ctx, msg, args...)
}

func (a *SlogAdapter) InfoContext(ctx context.Context, msg string, args ...any) {
	a.logger.InfoContext(ctx, msg, args...)
}

func (a *SlogAdapter) WarnContext(ctx context.Context, msg string, args ...any) {
	a.logger.WarnContext(ctx, msg, args...)
}

func (a *SlogAdapter) ErrorContext(ctx context.Context, msg string, args ...any) {
	a.logger.ErrorContext(ctx, msg, args...)
}

func (a *SlogAdapter) Fatal(msg string, args ...any) {
	a.logger.Error(msg, args...)
	os.Exit(1)
}

func (a *SlogAdapter) Fatalf(format string, args ...any) {
	a.logger.Error(fmt.Sprintf(format, args...))
	os.Exit(1)
}

func (a *SlogAdapter) FatalContext(ctx context.Context, msg string, args ...any) {
	a.logger.ErrorContext(ctx, msg, args...)
	os.Exit(1)
}

func (a *SlogAdapter) FatalfContext(ctx context.Context, format string, args ...any) {
	a.logger.ErrorContext(ctx, fmt.Sprintf(format, args...))
	os.Exit(1)
}

func (a *SlogAdapter) With(args ...any) ports.Logger {
	return &SlogAdapter{logger: a.logger.With(args...)}
}

func (a *SlogAdapter) WithGroup(name string) ports.Logger {
	return &SlogAdapter{logger: a.logger.WithGroup(name)}
}
