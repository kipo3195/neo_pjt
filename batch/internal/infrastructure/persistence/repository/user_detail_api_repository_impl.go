package repository

import (
	"batch/internal/consts"
	"batch/internal/domain/userDetail/repository"
	"bytes"
	"context"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

type userDetailApiRepositoryImpl struct {
	serverUrl string
}

func NewUserDetailApiRepository(serverUrl string) repository.UserDetailApiRepository {
	return &userDetailApiRepositoryImpl{
		serverUrl: serverUrl,
	}
}

func (r *userDetailApiRepositoryImpl) SendJsonToUser(ctx context.Context, fileName string, zipData []byte, orgCode string) error {
	url := "http://" + r.serverUrl + "/user/server/v1/detail/batch" // 상대 서버 endpoint

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// ---- 파일 파트 생성 ----
	part, err := writer.CreateFormFile("user_detail_file", fileName+".zip")
	if err != nil {
		return err
	}

	if _, err := part.Write(zipData); err != nil {
		return err
	}

	// ---- 추가 메타데이터 (선택) ----
	// 수신 쪽에서는 c.Request.FormValue(키워드)로 꺼내씀.
	_ = writer.WriteField("org_code", orgCode)

	// multipart 종료
	if err := writer.Close(); err != nil {
		return err
	}

	// ---- HTTP 요청 생성 ----
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		body,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	// 서버 키 인증 용
	// req.Header.Set("Authorization", "Bearer "+r.apiToken)

	// ---- 전송 ----
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseCode := resp.StatusCode
	responseBody, _ := io.ReadAll(resp.Body)
	log.Println("[SendJsonToUser] response code : ", responseCode)
	log.Println("[SendJsonToUser] response body : ", string(responseBody))

	if responseCode != http.StatusOK {
		return consts.ErrUserDetailFileSendError
	}

	return nil
}
