package router

import (
	"net/http"
	"notificator/internal/adapter/http/handler"
	"notificator/internal/adapter/http/middleware"
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

	//router.Use(middleware.LoggingMiddleware(logger))
	// 이건 gin 라이브러리와는 다른 gorilla/mux 라이브러리의 특징인데,
	// router.Use()로 등록된 미들웨어는 요청된 경로가 라우터에 등록된 경로와 일치(Match)할 때만 실행되며, 등록 되지 않은 endpoint는 미들웨어 체인을 호출하지 않는다.
	// 그러므로 전역적으로 설정하려면 조금 다르게 처리해야함.
	// mux.Router 자체도 결국 http.Handler. 라우터에게 일을 시키기 전에 밖에서 한 번 더 감싸버리면 404고 뭐고 무조건 로깅을 타게 됩니다.
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
