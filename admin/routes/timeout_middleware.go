package routes

import (
	"context"
	"net/http"
	"time"
)

func TimeoutMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 이 코드는 5초후에 채널이 ctx.Done()채널이 닫히면서 취소신호를 보냄.
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		// 그러므로 5초이상이 걸릴법한 비즈니스 로직에서 아래 주석과 같은 로직을 작성하여 처리 할 수 있음.

		// select {
		// case <-ctx.Done():
		// 	// 요청 취소 또는 타임아웃 발생
		// 	return ctx.Err()
		// default:
		// 계속 작업
		//}

		defer cancel()

		// r.WithContext(ctx)로 새 컨텍스트 넣기
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
