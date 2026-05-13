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

		parentCtx := c.Request.Context()
		ctx, cancel := context.WithTimeout(parentCtx, 10*time.Second)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		finished := make(chan struct{}, 1)

		go func() {
			c.Next()
			finished <- struct{}{}
		}()

		select {
		case <-finished:
			if ctx.Err() == context.DeadlineExceeded && !c.Writer.Written() {
				response.SendError(c, commonConsts.GATEWAY_TIMEOUT, commonConsts.ERROR, commonConsts.E_504, commonConsts.E_504_MSG)
				return
			}
			log.Println("request finished before timeout")
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				log.Println("request timeout")
				if c.Writer.Written() {
					return
				}
				response.SendError(c, commonConsts.GATEWAY_TIMEOUT, commonConsts.ERROR, commonConsts.E_504, commonConsts.E_504_MSG)
				return
			}
		}
	}
}
