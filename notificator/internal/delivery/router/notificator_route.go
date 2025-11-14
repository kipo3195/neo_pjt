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

func NewNotificatorRouter(serviceName string) notificatorRouter {
	router := mux.NewRouter()
	sub := router.PathPrefix("/" + serviceName).Subrouter()
	return notificatorRouter{
		R: sub,
	}
}

func (r *notificatorRouter) SetNotificatorServiceRoutes(handler *handler.NotificatorServiceHandler) {

	r.R.HandleFunc("/ws/connect", handler.NotificatorConnect)

}
