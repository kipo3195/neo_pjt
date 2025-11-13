package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	NATS     = "NATS"
	KAFKA    = "KAFKA"
	RABBITMQ = "RABBITMQ"
)

type ServerConfig struct {
	dbConfig *DBConfig
	mbConfig *MessageBrokerConfig
}

type DBConfig struct {
	Host     string
	Port     string
	Id       string
	Pw       string
	Database string
}

type MessageBrokerConfig struct {
	Mb string
}

func isLocal() bool {
	// Getenv는 로컬(gcp)과 운영(k8s)환경에서 다르게 동작함.
	// k8s에서는 이미 환경변수로 주입되어 있기 때문에 godotenv.Load()를 하지않아도 os.Getenv("값")을 호출해서 접근 할 수 있는 반면
	// 로컬 환경에서는 godotenv.Load()해야 .env 파일을 읽고 환경변수에 등록하기 때문에 이후에나 os.Getenv("값")에 호출 할 수 있게됨.
	return os.Getenv("ENV") == "" // 로컬일때는 환경변수에 등록되어 있지않으므로 .env에 ENV가 있어도 빈값을 반환함 그러므로 return true
}

func NewServerConfig() *ServerConfig {
	if isLocal() {
		godotenv.Load()
	}

	// DB 설정
	dbConfig := initDBConfig()
	mbConfig := initMBConfig()

	return &ServerConfig{
		dbConfig: dbConfig,
		mbConfig: mbConfig,
	}
}

/* DATABASE */
func initDBConfig() *DBConfig {
	host := os.Getenv("DB_HOST")
	id := os.Getenv("DB_ID")
	pw := os.Getenv("DB_PW")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_DATABASE")

	return &DBConfig{
		Host: host, Id: id, Pw: pw, Port: port, Database: database}
}

func ConnectDatabase(sfg *ServerConfig) *gorm.DB {

	dsn := sfg.dbConfig.Id + ":" + sfg.dbConfig.Pw + "@tcp(" + sfg.dbConfig.Host + ":" + sfg.dbConfig.Port + ")/" + sfg.dbConfig.Database + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}) // MYSQL

	if err != nil {
		log.Fatal("Failed to connect to database!")
	}

	log.Println("Message Database Connected !")
	return db
}

/* MESSAGE BROKER */
func initMBConfig() *MessageBrokerConfig {
	mb := os.Getenv("MB")

	return &MessageBrokerConfig{
		Mb: mb,
	}
}

func ConnectMessageBroker(sfg *ServerConfig) *nats.Conn {
	//broker.Broker {

	// 메시지 브로커 분기처리
	// switch sfg.mbConfig.Mb {
	// case NATS:

	connector, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Println("Failed to connect to NATS:", err)
		return nil
	}

	return connector
	// return &broker.NatsBroker{
	// 	Connector: connector,
	// }
	// case KAFKA:
	// 	log.Println("kafka is not available.")
	// case RABBITMQ:
	// 	log.Println("RabbitMQ is not available.")
	// }
	//	return nil
}
