package main

import (
	"context"
	"log"
	"net/http"
	"notificator/internal/di"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {

	log.Println("notificator 서비스 배포 테스트 1")

	modules, err := di.InitApp()

	if err != nil {
		// "DB 연결 실패", "설정 파일 누락" 등 구체적인 원인을 출력하고 종료
		log.Fatalf("Notificator service init error: %v", err)
	}

	go func() {
		log.Println("Notificator service is running on :8082")
		if err := modules.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Notificator service listen: %s\n", err)
		}
	}()

	// 세션당 메모리 사용량 측정용
	// go measure()

	// 시스템 시그널 대기 (SIGINT, SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Notificator service ...")

	// Graceful Shutdown 실행
	// HTTP 서버 먼저 종료 (새로운 요청 차단)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := modules.Server.Shutdown(ctx); err != nil {
		log.Fatal("Notificator service shutdown:", err)
	}

	// 5. 비동기 워커 풀 종료 (남은 작업 처리)
	modules.ChatModule.Cleanup()
	// 필요하다면 다른 모듈의 Cleanup도 호출
	// modules.NoteModule.Cleanup()

	log.Println("Notificator service exiting.")
}

func measure() {
	time.Sleep(5 * time.Second) // 서버 안정화

	before := printMem("before")

	log.Println("👉 3000 connections 붙이세요")
	time.Sleep(40 * time.Second)
	after1000 := printMem("after 3000")

	log.Println("===== RESULT =====")

	log.Printf("0→3000   : start=%d KB, end=%d KB, delta=%d KB, per=%d KB",
		before.HeapKB,
		after1000.HeapKB,
		after1000.HeapKB-before.HeapKB,
		(after1000.HeapKB-before.HeapKB)/3000,
	)

	// log.Printf("1000→2000: start=%d KB, end=%d KB, delta=%d KB, per=%d KB 🔥",
	// 	after1000.HeapKB,
	// 	after2000.HeapKB,
	// 	after2000.HeapKB-after1000.HeapKB,
	// 	(after2000.HeapKB-after1000.HeapKB)/1000,
	// )

	// log.Printf("2000→3000: start=%d KB, end=%d KB, delta=%d KB, per=%d KB",
	// 	after2000.HeapKB,
	// 	after3000.HeapKB,
	// 	after3000.HeapKB-after2000.HeapKB,
	// 	(after3000.HeapKB-after2000.HeapKB)/1000,
	// )
}

type Mem struct {
	HeapKB uint64
}

func printMem(tag string) Mem {
	runtime.GC()
	runtime.GC()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	heapKB := m.HeapAlloc / 1024

	log.Printf("[%s] HeapAlloc = %d KB", tag, heapKB)

	return Mem{
		HeapKB: heapKB,
	}
}
