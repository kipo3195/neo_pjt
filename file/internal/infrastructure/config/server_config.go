package config

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	// OCI SDK 필수 패키지
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/objectstorage"
)

type ServerConfig struct {
	Domain                string
	dbConfig              *DBConfig
	AutoMigrate           bool
	TokenConfig           TokenHashConfig
	OracleStorageConfig   OracleStorageConfig
	MessageFileGrpcConfig GrpcConfig
	BatchFileGrpcConfig   GrpcConfig
}

type OracleStorageConfig struct {
	Namespace  string
	BucketName string
	OciClient  objectstorage.ObjectStorageClient
}

type TokenHashConfig struct {
	AccessTokenHash  string
	RefreshTokenHash string
}

type DBConfig struct {
	Host     string
	Port     string
	Id       string
	Pw       string
	Database string
}

type GrpcConfig struct {
	Host string
	Port string
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

	autoMigrate := initAutoMigrate()

	tokenConfig := initTokenHash()

	oracleStorageConfig, err := initOracleStorageConfig()

	batchFileGrpcConfig := initBatchFileGrpcConfig()

	messageFileGrpcConfig := initMessageFileGrpcConfig()

	if err != nil {
		log.Panic(err)
	}

	return &ServerConfig{
		dbConfig:              dbConfig,
		AutoMigrate:           autoMigrate,
		TokenConfig:           tokenConfig,
		OracleStorageConfig:   oracleStorageConfig,
		BatchFileGrpcConfig:   batchFileGrpcConfig,
		MessageFileGrpcConfig: messageFileGrpcConfig,
	}
}

func initOracleStorageConfig() (OracleStorageConfig, error) {

	client, err := initOCIFromEnv()

	if err != nil {
		return OracleStorageConfig{}, err
	}

	namespace := os.Getenv("OCI_NAMESPACE")
	bucketName := os.Getenv("OCI_BUCKET_NAME")

	return OracleStorageConfig{
		OciClient:  client,
		Namespace:  namespace,
		BucketName: bucketName,
	}, nil

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

/* oracle cloud */
func initOCIFromEnv() (objectstorage.ObjectStorageClient, error) {
	// 환경 변수에서 직접 읽어서 설정 구성
	tenancy := os.Getenv("OCI_TENANCY_OCID")
	user := os.Getenv("OCI_USER_OCID")
	region := os.Getenv("OCI_REGION")
	fingerprint := os.Getenv("OCI_FINGERPRINT")

	privateKey := LoadPrivateKeyString()

	provider := common.NewRawConfigurationProvider(
		tenancy,
		user,
		region,
		fingerprint,
		privateKey,
		nil, // Passphrase가 없다면 nil
	)

	return objectstorage.NewObjectStorageClientWithConfigurationProvider(provider)
}

func LoadPrivateKeyString() string {
	keyPath := "/run/secrets/jwt_private_key"

	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		log.Fatalf("failed to read private key: %v", err)
	}

	return string(keyBytes)
}

func initTokenHash() TokenHashConfig {
	accessTokenHash := os.Getenv("ACCESS_TOKEN_HASH")
	refreshTokenHash := os.Getenv("REFRESH_TOKEN_HASH")
	return TokenHashConfig{
		AccessTokenHash:  accessTokenHash,
		RefreshTokenHash: refreshTokenHash,
	}
}

func ConnectDatabase(sfg *ServerConfig) (*gorm.DB, error) {

	//log.Println("env 읽음 " + sfg.dbConfig.Id + " : " + sfg.dbConfig.Pw + " : " + sfg.dbConfig.Host + " : " + sfg.dbConfig.Port + " : " + sfg.dbConfig.Database)
	dsn := sfg.dbConfig.Id + ":" + sfg.dbConfig.Pw + "@tcp(" + sfg.dbConfig.Host + ":" + sfg.dbConfig.Port + ")/" + sfg.dbConfig.Database + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}) // MYSQL

	if err != nil {
		return nil, err
	}

	log.Println("File Database Connected !")
	return db, nil
}

/* API Config Area*/

func initAutoMigrate() bool {
	flag := os.Getenv("AUTO_MIGRATE")
	if flag != "" && flag == "true" {
		return true
	} else {
		return false
	}
}

func ConnectCacheDataBase(sfg *ServerConfig) (*redis.ClusterClient, error) { // 반환 타입 변경

	// NewClusterClient는 바로 panic을 실행 시키지만,
	// Redis 서버 주소가 틀렸거나 서버가 꺼져 있는 경우 redis.ClusterClient 는 잘 만들어지므로 정상적으로 연결이 되었는지 검증 하는 걸 통해 error를 반환할 수 있도록 수정.

	if sfg == nil {
		return nil, errors.New("server config is nil")
	}

	// 추후 아래 설정을 sfg에서 불러오도록 수정필요.
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			"140.245.73.74:7001",
			"140.245.73.74:7002",
			"140.245.73.74:7003",
			"140.245.73.74:7004",
			"140.245.73.74:7005",
			"140.245.73.74:7006",
		},
	})

	// 2. 런타임 연결 확인 (패닉은 안 나지만 연결은 안 됐을 수 있음)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis connection error: %w", err)
	}

	return client, nil
}

func GetMessageFileLis() (net.Listener, error) {
	port := os.Getenv("MESSAGE_FILE_GRPC_PORT")
	return net.Listen("tcp", ":"+port)
}

func initMessageFileGrpcConfig() GrpcConfig {
	host := os.Getenv("MESSAGE_FILE_GRPC_HOST")
	port := os.Getenv("MESSAGE_FILE_GRPC_PORT")
	return GrpcConfig{
		Host: host,
		Port: port,
	}
}

func GetBatchFileLis() (net.Listener, error) {
	port := os.Getenv("BATCH_FILE_GRPC_PORT")
	return net.Listen("tcp", ":"+port)
}

func initBatchFileGrpcConfig() GrpcConfig {
	host := os.Getenv("BATCH_FILE_GRPC_HOST")
	port := os.Getenv("BATCH_FILE_GRPC_PORT")
	return GrpcConfig{
		Host: host,
		Port: port,
	}
}
