package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	ServerIP string
}

func isLocal() bool {
	// Getenv는 로컬(gcp)과 운영(k8s)환경에서 다르게 동작함.
	// k8s에서는 이미 환경변수로 주입되어 있기 때문에 godotenv.Load()를 하지않아도 os.Getenv("값")을 호출해서 접근 할 수 있는 반면
	// 로컬 환경에서는 godotenv.Load()해야 .env 파일을 읽고 환경변수에 등록하기 때문에 이후에나 os.Getenv("값")에 호출 할 수 있게됨.
	return os.Getenv("ENV") == "" // 로컬일때는 환경변수에 등록되어 있지않으므로 .env에 ENV가 있어도 빈값을 반환함 그러므로 return true
}

func NewServerConfig() *ServerConfig {

	if isLocal() {
		fmt.Println("11")
		godotenv.Load()
	}

	targetServerIP := os.Getenv("TARGET_SERVER_IP")

	fmt.Println("targetServerIP :", targetServerIP)

	return &ServerConfig{
		ServerIP: targetServerIP,
	}
}
