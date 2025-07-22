package config

import (
	"core/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServerConfig struct {
	dbConfig  *DBConfig
	ApiConfig *ApiConfig
}

type DBConfig struct {
	Host     string
	Port     string
	Id       string
	Pw       string
	Database string
}

type ApiConfig struct {
	NeoHeaderPrefix string
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

	// API 요청 설정
	apiConfig := initApiConfig()

	return &ServerConfig{
		dbConfig:  dbConfig,
		ApiConfig: apiConfig,
	}
}

/* DB Config Area*/

func initDBConfig() *DBConfig {
	host := os.Getenv("HOST")
	id := os.Getenv("ID")
	pw := os.Getenv("PW")
	port := os.Getenv("PORT")
	database := os.Getenv("DATABASE")

	return &DBConfig{
		Host:     host,
		Id:       id,
		Pw:       pw,
		Port:     port,
		Database: database,
	}
}

func ConnectDatabase(sfg *ServerConfig) *gorm.DB {

	//log.Println("env 읽음 " + sfg.dbConfig.Id + " : " + sfg.dbConfig.Pw + " : " + sfg.dbConfig.Host + " : " + sfg.dbConfig.Port + " : " + sfg.dbConfig.Database)
	dsn := sfg.dbConfig.Id + ":" + sfg.dbConfig.Pw + "@tcp(" + sfg.dbConfig.Host + ":" + sfg.dbConfig.Port + ")/" + sfg.dbConfig.Database + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}) // MYSQL

	if err != nil {
		log.Fatal("Failed to connect to database!")
	}

	/* 마이그레이션 */
	db.AutoMigrate(&models.AppValidation{})
	db.AutoMigrate(&models.WorksList{})

	log.Println("Core Database Connected !")
	return db
}

/* API Config Area*/

func initApiConfig() *ApiConfig {

	neoHeaderPrefix := os.Getenv("NEO_HEADER_PREFIX")

	return &ApiConfig{
		NeoHeaderPrefix: neoHeaderPrefix,
	}

}
