package response

import (
	"bytes"
	"fmt"
	"io"
	"message/pkg/consts"
	"message/pkg/dto"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// pkg/response: HTTP 응답 관련 로직 (추천)

// 순수한 HTTP 응답 처리
// 프레임워크(gin) 의존적
// 도메인에서 재사용 가능

func SendError(c *gin.Context, status int, result string, code string, msg string) {

	res := dto.ResponseDTO[dto.ErrorDataDTO]{ // 제네릭 타입 명시 - ResponseDTO의 DATA 'T'에 들어갈 타입을 말함.
		Result: result, // error, fail
		Data: dto.ErrorDataDTO{
			Code:    code,
			Message: msg,
		},
	}
	c.AbortWithStatusJSON(status, res)
	// 이미 응답을 보냈는데 뒤늦게 c.Next() 고루틴이 끝나면서 또 다른 응답을 보낼 수 없게 Abort
}

func SendSuccess[T any](c *gin.Context, t T) {
	res := dto.ResponseDTO[T]{ // 제네릭 타입 명시 - success는 어떤 DTO라도 들어갈 수 있으므로 any
		Result: consts.SUCCESS,
		Data:   t,
	}
	c.AbortWithStatusJSON(200, res) // 200 고정
	// 이미 응답을 보냈는데 뒤늦게 c.Next() 고루틴이 끝나면서 또 다른 응답을 보낼 수 없게 Abort
}

// 바이트 배열 전송 (메모리에 있는 파일을 전송)
func SendFileStream(c *gin.Context, data []byte, filename string, contentType string) {
	if contentType == "" {
		ext := filepath.Ext(filename)
		contentType = mime.TypeByExtension(ext)
		if contentType == "" {
			contentType = "application/octet-stream"
		}
	}

	c.Header("Content-Description", "File Stream")
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%q", filename))
	c.Header("Content-Type", contentType)
	c.Header("Cache-Control", "public, max-age=86400")

	c.Status(http.StatusOK)
	reader := bytes.NewReader(data)
	_, _ = io.Copy(c.Writer, reader)

}
