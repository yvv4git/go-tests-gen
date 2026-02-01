package ports

import "context"

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)

	DebugContext(ctx context.Context, msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)

	Fatal(msg string, args ...any)
	Fatalf(format string, args ...any)
	FatalContext(ctx context.Context, msg string, args ...any)
	FatalfContext(ctx context.Context, format string, args ...any)

	With(args ...any) Logger
	WithGroup(name string) Logger
}
