package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"test/internal/di"
)

func main() {

	log.Println("Test service init!")

	// 서버 및 모듈 초기화
	modules, err := di.InitApp()
	if err != nil {
		log.Println("Test service init error!")
		return
	}

	go func() {
		log.Println("Test service is running..")
		if err := modules.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Test service listen: %s\n", err)
		}
	}()

	// 시스템 시그널 대기 (SIGINT, SIGTERM) --------------------------------- 2
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// quit 채널에 데이터가 들어오기 전까지는 아래의 Shutdown 로직으로 넘어가지 않고 메인 함수가 계속 살아있게 됩니다. --------------------------------- 3
	<-quit

	// 시스템 종료신호가 들어왔을때 quit채널에 신호를 넣기 때문에 signal.Notify(quit...) <-quit가 풀리면서 하위 로직 수행 ----------------------- 4
	log.Println("Shutdown Test service ...")
}
