package config

import (
	"log"
	"os"

	"github.com/nats-io/nats.go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServerConfig struct {
	dbConfig    *DBConfig
	AutoMigrate bool
	Domain      string
}

type DBConfig struct {
	Host     string
	Port     string
	Id       string
	Pw       string
	Database string
}

func NewServerConfig() *ServerConfig {

	log.Println("Init admin ServerConfig !")
	// DB 설정을 읽어야함.
	dbConfig := &DBConfig{Host: "127.0.0.1", Id: "neo", Pw: "neo", Port: "3306", Database: "admin"}

	autoMigrate := initAutoMigrate()

	domain := initDomain()

	return &ServerConfig{
		dbConfig:    dbConfig,
		AutoMigrate: autoMigrate,
		Domain:      domain,
	}
}

func ConnectMessageBroker(sfg *ServerConfig) *nats.Conn {

	// 메시지 브로커 분기처리
	// switch sfg.mbConfig.Mb {
	// case NATS:
	nc, err := nats.Connect(nats.DefaultURL)

	if err != nil {
		log.Println("Failed to connect to NATS:", err)
		return nil
	}
	return nc
	// 	return &broker.NatsBroker{
	// 		Nc:        nc,
	// 		ChatUsers: make(map[string]*broker.ChatUser),
	// 	}
	// case KAFKA:
	// 	log.Println("kafka is not available.")
	// case RABBITMQ:
	// 	log.Println("RabbitMQ is not available.")
	// }
	// return nil
}

func ConnectDatabase(sfg *ServerConfig) *gorm.DB {

	dsn := sfg.dbConfig.Id + ":" + sfg.dbConfig.Pw + "@tcp(" + sfg.dbConfig.Host + ":" + sfg.dbConfig.Port + ")/" + sfg.dbConfig.Database + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}) // MYSQL

	if err != nil {
		log.Fatal("Failed to connect to database!")
	}

	log.Println("Admin Database Connected !")
	return db
}

func initAutoMigrate() bool {
	flag := os.Getenv("AUTO_MIGRATE")
	if flag != "" && flag == "true" {
		return true
	} else {
		return false
	}
}

func initDomain() string {
	return os.Getenv("DOMAIN")
}
