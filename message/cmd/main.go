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

	log.Println("자동화 테스트 점검용 로깅 7")

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

	/* Graceful ShutDown 로직 */

	// HTTP 서버 실행 (비동기)
	go func() {
		log.Println("Message service is running..")
		// ListenAndServe는 서버가 종료될때까지 대기상태로 블로킹 단, 별도 고루틴으로 실행시켰으므로 아래 로직으로 내려감  ---------------------- 1
		if err := modules.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// 모든 요청이 처리되면 (ctx를 생성한 최대 10초 까지만) 비동기로 돌던 ListenAndServe()가 http.ErrServerClosed 에러를 반환하며 종료 -------------------------- 6
			log.Fatalf("Message service listen: %s\n", err)
		}
	}()

	go func() {
		log.Println("[Batch - Message] gRPC is running..")
		// Serve는 서버가 종료될때까지 대기상태로 블로킹 단, 별도 고루틴으로 실행시켰으므로 아래 로직으로 내려감  ---------------------- 2
		if err := modules.BatchMessageGrpcServer.Serve(modules.BatchMessageListener); err != nil {
			// for select문에서 모든 처리가 완료되어 stopped 데이터가 들어오거나, time out 됬을때 ---------------------- 8-1
			log.Fatalf("Batch - Message gRPC serve error: %v", err)
		}
	}()

	// 시스템 시그널 대기 (SIGINT, SIGTERM) --------------------------------- 2
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// quit 채널에 데이터가 들어오기 전까지는 아래의 Shutdown 로직으로 넘어가지 않고 메인 함수가 계속 살아있게 됩니다. --------------------------------- 3
	<-quit

	// 시스템 종료신호가 들어왔을때 quit채널에 신호를 넣기 때문에 signal.Notify(quit...) <-quit가 풀리면서 하위 로직 수행 ----------------------- 4
	log.Println("Shutdown Message service ...")

	// Graceful Shutdown 실행
	// HTTP 서버 먼저 종료 (새로운 요청 차단)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := modules.Server.Shutdown(ctx); err != nil { // Shutdown이 호출되는 즉시, 서버는 더 이상 새로운 HTTP 연결을 받지 않음. -------------------------- 5
		log.Fatal("Message service shutdown:", err)
	}

	// gRPC GracefulStop 실행
	// GracefulStop은 처리가 완료될때까지 무기한 대기 별도의 자체 타임아웃 없음 그러므로, 별도의 고루틴으로 실행하고 타임아웃을 걸어서 일정 시간 동안만 유지.
	batchMessageStopped := make(chan struct{})
	go func() {
		modules.BatchMessageGrpcServer.GracefulStop() // GracefulStop 호출되는 즉시, 서버는 더 이상 새로운 gRPC 연결을 받지 않음. -------------------------- 7
		close(batchMessageStopped)
	}()

	// 서버가 모두 끝날 때까지 대기 (타임아웃 적용)
	timeout := time.After(5 * time.Second)
	select {
	case <-batchMessageStopped:
		log.Println("Batch - Message gRPC server stopped")
	case <-timeout:
		log.Println("gRPC stop timeout, forcing stop")
		modules.BatchMessageGrpcServer.Stop()
	}

	// 자원 해제 호출  -------------------------------------- 7
	// ctx의 시간 (최대 10초 동안) 최대한 남은 일을 처리하려고 시간을 줘놓고 cleanup을 먼저 수행해버리면 안됨.
	modules.Cleanup()

	log.Println("Message service exiting.")
}
