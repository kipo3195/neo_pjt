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
	AuthenticateToken(token string) (bool, error)
	HandleAuth(conn *websocket.Conn, data map[string]interface{}) (bool, error)
}

func NewAuthUsecase(repo repositories.AuthRepository) AuthUsecase {
	return &authUsecase{repo: repo}
}

func (uc *authUsecase) AuthenticateToken(token string) (bool, error) {
	// 토큰 검증 로직

	if token == "" {
		return false, nil
	}

	return true, nil
}

func (uc *authUsecase) HandleAuth(conn *websocket.Conn, data map[string]interface{}) (bool, error) {
	token := data["token"].(string)

	_, err := uc.AuthenticateToken(token)
	if err != nil {
		// defer conn.Close() 호출 되므로 연결 종료.
		return false, errors.New("auth_failed")
	}
	return true, nil

}
