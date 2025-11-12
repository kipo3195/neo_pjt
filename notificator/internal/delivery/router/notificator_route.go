package router

import (
	"notificator/internal/delivery/handler"

	"github.com/gorilla/mux"
)

type notificatorRouter struct {
	R *mux.Router
}

type NotificatorRouter interface {
	SetNotificatorServiceRoutes(handler *handler.NotificatorServiceHandler)
}

func (r *notificatorRouter) SetNotificatorServiceRoutes(handler *handler.NotificatorServiceHandler) {

}

func NewNotificatorRouter(serviceName string) notificatorRouter {
	router := mux.NewRouter()
	return notificatorRouter{
		R: router,
	}
}
