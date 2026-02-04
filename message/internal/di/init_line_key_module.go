package di

import (
	"log"
	"message/internal/adapter/http/handler"
	"message/internal/application/usecase"
	"message/internal/infrastructure/persistence/repository"
	"message/pkg/util"

	"gorm.io/gorm"
)

type LineKeyModule struct {
	Handler *handler.LineKeyHandler
	Usecase usecase.LineKeyUsecase
}

func InitLineKeyModule(db *gorm.DB) *LineKeyModule {

	// ULID Generator는 서버 시작 시 한 번만 초기화
	ulidGen, err := util.NewULIDGenerator()
	if err != nil {
		panic("failed to init ULID generator: " + err.Error())
	}

	log.Println("ULID Generate Init !")

	repository := repository.NewLineKeyRepository(db)
	usecase := usecase.NewLineKeyUsecase(repository, ulidGen)
	handler := handler.NewLineKeyHandler(usecase)

	return &LineKeyModule{
		Handler: handler,
		Usecase: usecase,
	}
}
