package router

import (
	"net/http"
	"notificator/internal/delivery/handler"
	"notificator/internal/delivery/middleware"
	"notificator/internal/domain/logger"
	"notificator/internal/infrastructure/config"

	"github.com/gorilla/mux"
)

type notificatorRouter struct {
	R           *mux.Router
	tokenConfig config.TokenHashConfig
	logger      logger.Logger
}

type NotificatorRouter interface {
	SetNotificatorServiceRoutes(handler *handler.NotificatorServiceHandler)
}

func NewNotificatorRouter(serviceName string, tokenConfig config.TokenHashConfig, logger logger.Logger) notificatorRouter {
	router := mux.NewRouter()
	sub := router.PathPrefix("/" + serviceName).Subrouter()
	return notificatorRouter{
		R:           sub,
		tokenConfig: tokenConfig,
		logger:      logger,
	}
}

func (r *notificatorRouter) SetNotificatorServiceRoutes(handler *handler.NotificatorServiceHandler) {

	client := r.R.PathPrefix("/ws").Subrouter()
	client.Use(middleware.LoggingMiddleware(r.logger))
	client.Use(middleware.AuthMiddleware(r.logger, r.tokenConfig))
	client.HandleFunc("/connect", handler.NotificatorConnect).Methods(http.MethodGet)

}
