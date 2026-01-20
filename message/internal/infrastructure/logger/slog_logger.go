package logger

import (
	"context"
	"log/slog"
	"message/internal/domain/logger"
	"os"
)

type slogLogger struct {
	handler *slog.Logger // slog 라이브러리 사용시
}

func (l *slogLogger) Info(ctx context.Context, msg string, args ...any) {
	l.handler.InfoContext(ctx, msg)
}

func (l *slogLogger) Error(ctx context.Context, msg string, args ...any) {
	l.handler.ErrorContext(ctx, msg)
}

func NewslogLogger() logger.Logger {
	h := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	return &slogLogger{
		handler: h,
	}
}
