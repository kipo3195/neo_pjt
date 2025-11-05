package routes

import (
	"message/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes(messageHandler *handlers.MessageHandler) *mux.Router {
	// mux.Router는 기본적으로 HTTP 요청 경로를 라우팅하는 데 사용됩니다

	r := mux.NewRouter()
	msgV1 := r.PathPrefix("/msg/v1").Subrouter()
	msgV1.HandleFunc("/connect", messageHandler.HandleWebSocket)
	return r
}
