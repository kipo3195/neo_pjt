package usecases

import (
	"errors"
	"message/repositories"

	"github.com/gorilla/websocket"
)

type authUsecase struct {
	repo repositories.AuthRepository
}
type AuthUsecase interface {
	HandleAuth(conn *websocket.Conn, data map[string]interface{}) (bool, error)
}

func NewAuthUsecase(repo repositories.AuthRepository) AuthUsecase {
	return &authUsecase{repo: repo}
}

func (r *authUsecase) HandleAuth(conn *websocket.Conn, data map[string]interface{}) (bool, error) {
	token := data["token"].(string)

	_, err := authenticateToken(token)
	if err != nil {
		// defer conn.Close() 호출 되므로 연결 종료.
		return false, errors.New("auth_failed")
	}
	return true, nil

}

func authenticateToken(token string) (bool, error) {
	// 토큰 검증 로직 TODO

	if token == "" {
		return false, nil
	}

	return true, nil
}
