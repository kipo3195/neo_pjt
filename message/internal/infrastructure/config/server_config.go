package config

import (
	"log"
	"message/internal/infrastructure/broker"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServerConfig struct {
	dbConfig    *DBConfig
	jwtConfig   *JWTConfig // jwt를 소문자로 정의했으므로 외부에서 접근할 수 없음 = 그래서 GetJWTConfig 메소드를 만들어서 외부에서 사용 할 수 있게함.
	AutoMigrate bool
	TokenConfig TokenHashConfig
	mbConfig    *MessageBrokerConfig
}

type TokenHashConfig struct {
	AccessTokenHash  string
	RefreshTokenHash string
}

type MessageBrokerConfig struct {
	Mb string
}

type DBConfig struct {
	Host     string
	Port     string
	Id       string
	Pw       string
	Database string
}

type JWTConfig struct {
	Key                string
	AccessExp          int
	RefressExp         int
	AppTokenExp        int
	AppRefreshTokenExp int
}

func isLocal() bool {
	// Getenv는 로컬과 운영(k8s)환경에서 다르게 동작함.
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
	// jwt 설정
	jwtConfig := initJwtConfig()

	autoMigrate := initAutoMigrate()

	tokenConfig := initTokenHash()

	return &ServerConfig{
		dbConfig:    dbConfig,
		jwtConfig:   jwtConfig,
		AutoMigrate: autoMigrate,
		TokenConfig: tokenConfig,
	}
}
func initDBConfig() *DBConfig {
	host := os.Getenv("HOST")
	id := os.Getenv("ID")
	pw := os.Getenv("PW")
	port := os.Getenv("PORT")
	database := os.Getenv("DATABASE")

	return &DBConfig{
		Host: host, Id: id, Pw: pw, Port: port, Database: database}
}

func initTokenHash() TokenHashConfig {
	accessTokenHash := os.Getenv("ACCESS_TOKEN_HASH")
	refreshTokenHash := os.Getenv("REFRESH_TOKEN_HASH")
	return TokenHashConfig{
		AccessTokenHash:  accessTokenHash,
		RefreshTokenHash: refreshTokenHash,
	}
}

// jwt 설정 조회
func initJwtConfig() *JWTConfig {
	// 이 값이 환경변수에 등록되어있지않으면 검증하는 쪽에서 signature에러 발생.
	key := os.Getenv("JWT_SECRET_KEY")
	accessExp, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXP_M"))
	if err != nil {
		log.Println("ACCESS_TOKEN_EXP_M is invalid. :", err)
	}
	refreshExp, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXP_D"))
	if err != nil {
		log.Println("REFRESH_TOKEN_EXP_D is invalid. :", err)
	}
	appTokenExp, err := strconv.Atoi(os.Getenv("APP_TOKEN_EXP_D"))
	if err != nil {
		log.Println("REFRESH_TOKEN_EXP_D is invalid. :", err)
	}
	appRefreshTokenExp, err := strconv.Atoi(os.Getenv("APP_REFRESH_TOKEN_EXP_D"))
	if err != nil {
		log.Println("REFRESH_TOKEN_EXP_D is invalid. :", err)
	}
	return &JWTConfig{Key: key, AccessExp: accessExp, RefressExp: refreshExp, AppTokenExp: appTokenExp, AppRefreshTokenExp: appRefreshTokenExp}
}

func (s *ServerConfig) GetJWTConfig() *JWTConfig {
	return s.jwtConfig
}

func ConnectDatabase(sfg *ServerConfig) *gorm.DB {

	//log.Println("env 읽음 " + sfg.dbConfig.Id + " : " + sfg.dbConfig.Pw + " : " + sfg.dbConfig.Host + " : " + sfg.dbConfig.Port + " : " + sfg.dbConfig.Database)
	dsn := sfg.dbConfig.Id + ":" + sfg.dbConfig.Pw + "@tcp(" + sfg.dbConfig.Host + ":" + sfg.dbConfig.Port + ")/" + sfg.dbConfig.Database + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}) // MYSQL

	if err != nil {
		log.Fatal("Failed to connect to database!")
	}

	//db.AutoMigrate(&sharedModels.ServiceUsers{})

	log.Println("Message Database Connected !")
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

func ConnectMessageBroker(sfg *ServerConfig) broker.Broker {

	// 메시지 브로커 분기처리
	switch sfg.mbConfig.Mb {
	case NATS:
		nc, err := nats.Connect(nats.DefaultURL)
		if err != nil {
			log.Println("Failed to connect to NATS:", err)
			return nil
		}
		return &broker.NatsBroker{
			Nc:        nc,
			ChatUsers: make(map[string]*broker.ChatUser),
		}
	case KAFKA:
		log.Println("kafka is not available.")
	case RABBITMQ:
		log.Println("RabbitMQ is not available.")
	}
	return nil
}

const (
	NATS     = "NATS"
	KAFKA    = "KAFKA"
	RABBITMQ = "RABBITMQ"
)
