package middleware

import (
	"context"
	"log"
	commonConsts "message/pkg/consts"
	"message/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
)

func TimeoutMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		// 현재 TimeoutMiddleware의 부모 격인 LoggerMiddleware의 context를 가져옴
		parentCtx := c.Request.Context()

		// 덧붙여 ctx 생성
		ctx, cancel := context.WithTimeout(parentCtx, 10*time.Second)
		// 10초가 지났을 때 (타임아웃)
		// 비즈니스 로직이 아직 돌고 있는데 10초가 딱 됩니다.
		// 타임아웃 시점(10초가 되는 순간)**에 Go 내부 타이머가 이미
		// ctx.Err()를 context.DeadlineExceeded로 채워버립니다. 즉, defer cancel()이 실행되기 전에 이미 하위 레이어들은 에러를 인지할 수 있는 상태가 됩니다.
		// select 문에서 case <-ctx.Done():이 즉시 실행되어 504 에러를 보냅니다.

		// defer cancel() "이미 취소된 거 다시 확인사살하고, 타이머 리소스를 시스템에 반환"하는 마무리 작업만 합니다. 에러 발생 시점에는 영향을 주지 않습니다.
		defer cancel()

		// 새 ctx를 주입
		c.Request = c.Request.WithContext(ctx)

		finished := make(chan struct{}, 1)

		go func() {
			// 모든 요청이 끝나면 여기가 먼저 호출이 되는듯?
			// c.Next()를 그냥 호출하면, 하위 핸들러가 10초를 넘게 써도 c.Next()가 끝날 때까지 미들웨어는 기다리게 됩니다.
			// 그러므로 별도의 고루틴에서 하위 핸들러 처리가 끝났을때 finished의 채널에 데이터를 넣어주는 처리를 통해 모든 handler 이벤트가 끝났음을 명시적으로 작성합니다.
			c.Next()
			// 이 시점이 비즈니스 로직 완료 후 호출되는 시점인듯
			finished <- struct{}{}
		}()

		select {
		case <-finished:
			log.Println("요청이 먼저 끝남")
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				log.Println("시간 초과로 인해 ctx 중지 호출이 먼저 옴 ")
				response.SendError(c, commonConsts.GATEWAY_TIMEOUT, commonConsts.ERROR, commonConsts.E_504, commonConsts.E_504_MSG)
				return
			}
		}
	}
}
