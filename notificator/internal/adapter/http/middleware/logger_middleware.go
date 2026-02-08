package middleware

import (
	"context"
	"net/http"
	"notificator/internal/adapter/http/middleware/wrapper"
	"notificator/internal/domain/logger"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func LoggingMiddleware(logger logger.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()
			traceID := uuid.New().String()

			// 1. TraceID를 표준 context에 주입
			ctx := context.WithValue(r.Context(), "trace_id", traceID)

			// 2. Gin Context의 Request를 새로운 context로 교체
			r = r.WithContext(ctx)

			rw := wrapper.NewResponseWriter(w)

			// defer를 통해 모든 핸들러가 끝난 후 로그 출력
			defer func() {

				// c.Next()를 거치며 미들웨어가 추가하더라도 그 미들웨어에서 업데이트한 데이터를 꺼내 오기 위함
				latestCtx := r.Context()

				// 성공이든 실패든 무조건 실행되는 '최종 결과 요약'
				// 상태 코드가 400, 500이면 실패인 걸 이미 status 필드가 말해주고 있습니다.
				logger.Info(latestCtx, "request_finished",
					"method", r.Method,
					"path", r.URL.Path,
					"status", rw.StatusCode, // 여기서 200인지 500인지 찍힘
					"latency", time.Since(start),
				)
			}()

			// 다음 미들웨어 또는 핸들러 실행
			next.ServeHTTP(w, r)
		})
	}
}
