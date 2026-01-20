package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		traceID := uuid.New().String()

		// TraceID를 Context에 주입하여 하위 레이어로 전달
		ctx := context.WithValue(r.Context(), "trace_id", traceID)

		// 실제 핸들러 수행 (defer를 통해 응답 후 로그 출력)
		defer func() {
			logger.InfoContext(ctx, "request_completed",
				"trace_id", traceID,
				"method", r.Method,
				"path", r.URL.Path,
				"latency", time.Since(start),
			)
		}()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
