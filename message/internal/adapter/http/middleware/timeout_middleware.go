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
		// Go 내부 타이머가 즉시 ctx.Done() 채널에 신호를 보냅니다. (이게 사실상 내부적인 cancel 호출입니다.)
		// select 문에서 case <-ctx.Done():이 즉시 실행되어 504 에러를 보냅니다.
		// 함수가 종료되면서 아래 defer cancel()이 한 번 더 확실하게 자원을 정리합니다.
		defer cancel()
		// 그럼 미들웨어를 필두로하여 ctx를 전달 받은 handler -> usecase -> repository 순으로 전파됨
		// 각 layer의 ctx 에러 처리 로직 발동.

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
