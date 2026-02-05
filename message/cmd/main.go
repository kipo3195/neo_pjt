package main

import (
	"context"
	"log"
	"message/internal/di"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// 서버 및 모듈 초기화
	modules, err := di.InitApp()

	if err != nil {
		// "DB 연결 실패", "설정 파일 누락" 등 구체적인 원인을 출력하고 종료
		log.Fatalf("Message service init error: %v", err)
	}

	// 서버 실행 전 데이터 로딩
	// 메모리에 데이터가 올라가야 하므로 ListenAndServe 이전에 실행
	ctxLoad, cancelLoad := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelLoad()

	log.Println("Message service loading initial data...")
	// 여기서 생성된 context를 loader에도 그대로 주입 (AS-IS 로직에서는 ctx 주입하지 않았었음..)
	if err := modules.DataLoader.LoadAllData(ctxLoad); err != nil {
		// 데이터 로딩 실패 시 서버를 띄우지 않고 종료하는 것이 안전하므로 종료 처리
		log.Fatalf("Message service data loading failed: %v", err)
	}

	// 서버 실행 (비동기)
	go func() {
		log.Println("Message service is running..")
		if err := modules.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Message service listen: %s\n", err)
		}
	}()

	// 시스템 시그널 대기 (SIGINT, SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Message service ...")

	// Graceful Shutdown 실행
	// HTTP 서버 먼저 종료 (새로운 요청 차단)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := modules.Server.Shutdown(ctx); err != nil {
		log.Fatal("Message service shutdown:", err)
	}

	// 비동기 워커 풀 종료 (남은 작업 처리)
	modules.Cleanup()

	// 필요하다면 다른 모듈의 Cleanup도 호출
	// modules.NoteModule.Cleanup()
	log.Println("Message service exiting.")
}
