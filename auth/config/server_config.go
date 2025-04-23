package config

import (
	"auth/models"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServerConfig struct {
	dbConfig *DBConfig
}

type DBConfig struct {
	Host     string
	Port     string
	Id       string
	Pw       string
	Database string
}

func NewServerConfig() *ServerConfig {
	fmt.Println("Init auth serverConfig")
	dbConfig := &DBConfig{Host: "127.0.0.1", Id: "neo", Pw: "neo", Port: "3306", Database: "auth"}

	return &ServerConfig{
		dbConfig: dbConfig,
	}
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
