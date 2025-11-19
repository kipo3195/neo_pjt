package response

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"notificator/pkg/consts"
	"notificator/pkg/dto"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// pkg/response: HTTP 응답 관련 로직 (추천)

// 순수한 HTTP 응답 처리
// 프레임워크(gin) 의존적
// 도메인에서 재사용 가능

func SendError(w http.ResponseWriter, status int, result string, code string, msg string) {
	// Gin의 dto 구조체 그대로 사용한다고 가정
	res := dto.ResponseDTO[dto.ErrorDataDTO]{
		Result: result, // "error", "fail" 등
		Data: dto.ErrorDataDTO{
			Code:    code,
			Message: msg,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// JSON 응답 출력
	json.NewEncoder(w).Encode(res)
}

func SendSuccess[T any](conn *websocket.Conn, t T) {
	res := dto.ResponseDTO[T]{ // 제네릭 타입 명시 - success는 어떤 DTO라도 들어갈 수 있으므로 any
		Result: consts.SUCCESS,
		Data:   t,
	}
	conn.WriteJSON(res)
}

func SendErrorWs(conn *websocket.Conn, status int, result string, code string, msg string) {
	res := dto.ResponseDTO[dto.ErrorDataDTO]{
		Result: result, // "error", "fail" 등
		Data: dto.ErrorDataDTO{
			Code:    code,
			Message: msg,
		},
	}
	conn.WriteJSON(res)
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
