package router

import (
	"notificator/internal/delivery/handler"
	"notificator/internal/delivery/middleware"
	"notificator/internal/infrastructure/config"

	"github.com/gorilla/mux"
)

type notificatorRouter struct {
	R           *mux.Router
	tokenConfig config.TokenHashConfig
}

type NotificatorRouter interface {
	SetNotificatorServiceRoutes(handler *handler.NotificatorServiceHandler)
}

func NewNotificatorRouter(serviceName string, tokenConfig config.TokenHashConfig) notificatorRouter {
	router := mux.NewRouter()
	sub := router.PathPrefix("/" + serviceName).Subrouter()
	return notificatorRouter{
		R:           sub,
		tokenConfig: tokenConfig,
	}
}

func (r *notificatorRouter) SetNotificatorServiceRoutes(handler *handler.NotificatorServiceHandler) {

	// REST API와는 다른 인증 middleware 구조
	r.R.HandleFunc("/ws/connect", middleware.AuthMiddleware(handler.NotificatorConnect, r.tokenConfig))
}
