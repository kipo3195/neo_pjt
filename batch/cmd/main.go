package main

import (
	"batch/internal/di"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	modules, err := di.InitApp()

	if err != nil {
		log.Fatalf("Batch service init error :%v", err)
	}

	// HTTP 서버 실행 (비동기)
	go func() {
		log.Println("Batch service is running ...")
		if err := modules.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Batch service listen: %s\n", err)
		}
	}()

	// 시스템 시그널 대기 (SIGINT, SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // 여기서 블로킹됨

	// Graceful Shutdown 실행
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := modules.Server.Shutdown(ctx); err != nil {
		log.Fatal("Batch service shutdown:", err)
	}

	modules.Cleanup()

	log.Println("Batch service exiting.")
}
