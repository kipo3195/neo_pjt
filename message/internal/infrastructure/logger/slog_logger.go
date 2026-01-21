package logger

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"time"
)

type slogLogger struct {
	handler slog.Handler
}

func (l *slogLogger) log(ctx context.Context, level slog.Level, msg string, args ...any) {
	// 현재 레벨이 활성화되어 있는지 확인
	if !l.handler.Enabled(ctx, level) {
		return
	}

	var pcs [1]uintptr
	// Skip 갯수:
	// [0]runtime.Callers -> [1]slogLogger.log -> [2]slogLogger.Info/Error -> [3]UseCase
	runtime.Callers(3, pcs[:])

	// 레코드 생성 (pcs[0]를 통해 호출지점 주입)
	r := slog.NewRecord(time.Now(), level, msg, pcs[0])

	// 컨텍스트에서 trace_id를 추출하여 로그 레코드에 추가
	if tid, ok := ctx.Value("trace_id").(string); ok {
		r.AddAttrs(slog.String("trace_id", tid))
	}

	// 외부에서 넘겨준 args(method, status, latency 등)를 레코드에 추가
	r.Add(args...)

	// 핸들러를 통해 최종 출력
	_ = l.handler.Handle(ctx, r)
}

func (l *slogLogger) Info(ctx context.Context, msg string, args ...any) {
	l.log(ctx, slog.LevelInfo, msg, args...)
}

func (l *slogLogger) Error(ctx context.Context, msg string, args ...any) {
	l.log(ctx, slog.LevelError, msg, args...)
}

func NewSlogLogger() *slogLogger {
	// JSON 핸들러를 생성합니다.
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	})
	return &slogLogger{handler: h}
}
