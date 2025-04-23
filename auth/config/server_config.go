package config

import (
	"auth/models"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServerConfig struct {
	dbConfig  *DBConfig
	jwtConfig *JWTConfig // jwt를 소문자로 정의했으므로 외부에서 접근할 수 없음 = 그래서 GetJWTConfig 메소드를 만들어서 외부에서 사용 할 수 있게함.
}

type DBConfig struct {
	Host     string
	Port     string
	Id       string
	Pw       string
	Database string
}

type JWTConfig struct {
	AccessExp  int
	RefressExp int
}

func NewServerConfig() *ServerConfig {
	fmt.Println("Init auth serverConfig")

	// DB 설정
	dbConfig := &DBConfig{Host: "127.0.0.1", Id: "neo", Pw: "neo", Port: "3306", Database: "auth"}

	// jwt 설정
	jwtConfig := &JWTConfig{AccessExp: 1, RefressExp: 30}

	return &ServerConfig{
		dbConfig:  dbConfig,
		jwtConfig: jwtConfig,
	}
}

func (s *ServerConfig) GetJWTConfig() *JWTConfig {
	return s.jwtConfig
}

func ConnectDatabase(sfg *ServerConfig) *gorm.DB {

	dsn := sfg.dbConfig.Id + ":" + sfg.dbConfig.Pw + "@tcp(" + sfg.dbConfig.Host + ":" + sfg.dbConfig.Port + ")/" + sfg.dbConfig.Database + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}) // MYSQL

	if err != nil {
		log.Fatal("Failed to connect to database!")
	}

	db.AutoMigrate(&models.AuthInfo{})

	fmt.Println("Auth Database Connected !")
	return db
}
