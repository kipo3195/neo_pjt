package config

import (
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServerConfig struct {
	Domain                string
	dbConfig              *DBConfig
	AutoMigrate           bool
	TokenConfig           TokenHashConfig
	OrgInfoBatchConfig    *BatchConfig
	UserDetailBatchConfig *BatchConfig
	ChatFileConfig        ChatFileConfig
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

type BatchConfig struct {
	Org              string         // works code
	BatchFlag        bool           // 배치 실행 여부
	Cron             string         // 실행 주기 데이터
	ExtendDBSyncFlag bool           // 외부 DB 연동 설정
	ExtendDBConfig   ExtendDBConfig // 외부 DB 연결 정보
}

type ChatFileConfig struct {
	FileServiceGrpcHost    string
	FileServiceGrpcPort    string
	MessageServiceGrpcHost string
	MessageServiceGrpcPort string
	ChatFileBatchConfig    *BatchConfig
}

// 공통 - 외부 DB 연동을 위한 정보
type ExtendDBConfig struct {
	Host     string
	Port     string
	Id       string
	Pw       string
	Database string
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

	domain := initDomain()

	orgInfoBatchConfig := initOrgInfoBatchConfig()

	userDetailBatchConfig := initUserDetailBatchConfig()

	chatFileConfig := initChatFileConfig()

	return &ServerConfig{
		dbConfig:              dbConfig,
		AutoMigrate:           autoMigrate,
		TokenConfig:           tokenConfig,
		OrgInfoBatchConfig:    &orgInfoBatchConfig,
		UserDetailBatchConfig: &userDetailBatchConfig,
		Domain:                domain,
		ChatFileConfig:        chatFileConfig,
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

func initUserDetailBatchConfig() BatchConfig {
	org := os.Getenv("USER_DETAIL")
	flag := os.Getenv("USER_DETAIL_CRON_BATCH_FLAG")
	cron := os.Getenv("USER_DETAIL_CRON")

	var batchFlag bool
	if flag == "true" {
		batchFlag = true
	} else {
		batchFlag = false
	}

	flag = os.Getenv("USER_DETAIL_EXTEND_DB_SYNC_FLAG")

	var extendDbSyncFlag bool
	if flag == "true" {
		extendDbSyncFlag = true
	} else {
		extendDbSyncFlag = false
	}

	extendDbHost := os.Getenv("USER_DETAIL_EXTEND_DB_HOST")
	extendDbPort := os.Getenv("USER_DETAIL_EXTEND_DB_PORT")
	extendDbId := os.Getenv("USER_DETAIL_EXTEND_DB_ID")
	extendDbPw := os.Getenv("USER_DETAIL_EXTEND_DB_PW")
	extendDbDatabase := os.Getenv("USER_DETAIL_EXTEND_DB_DATABASE")

	extendDBConfig := ExtendDBConfig{
		Host:     extendDbHost,
		Port:     extendDbPort,
		Id:       extendDbId,
		Pw:       extendDbPw,
		Database: extendDbDatabase,
	}

	return BatchConfig{
		Org:              org,
		BatchFlag:        batchFlag,
		Cron:             cron,
		ExtendDBSyncFlag: extendDbSyncFlag,
		ExtendDBConfig:   extendDBConfig,
	}
}

func initOrgInfoBatchConfig() BatchConfig {

	org := os.Getenv("ORG_INFO")
	flag := os.Getenv("ORG_INFO_CRON_BATCH_FLAG")
	cron := os.Getenv("ORG_INFO_CRON")

	var batchFlag bool
	if flag == "true" {
		batchFlag = true
	} else {
		batchFlag = false
	}

	flag = os.Getenv("ORG_INFO_EXTEND_DB_SYNC_FLAG")

	var extendDbSyncFlag bool
	if flag == "true" {
		extendDbSyncFlag = true
	} else {
		extendDbSyncFlag = false
	}

	extendDbHost := os.Getenv("ORG_INFO_EXTEND_DB_HOST")
	extendDbPort := os.Getenv("ORG_INFO_EXTEND_DB_PORT")
	extendDbId := os.Getenv("ORG_INFO_EXTEND_DB_ID")
	extendDbPw := os.Getenv("ORG_INFO_EXTEND_DB_PW")
	extendDbDatabase := os.Getenv("ORG_INFO_EXTEND_DB_DATABASE")

	extendDBConfig := ExtendDBConfig{
		Host:     extendDbHost,
		Port:     extendDbPort,
		Id:       extendDbId,
		Pw:       extendDbPw,
		Database: extendDbDatabase,
	}

	return BatchConfig{
		Org:              org,
		BatchFlag:        batchFlag,
		Cron:             cron,
		ExtendDBSyncFlag: extendDbSyncFlag,
		ExtendDBConfig:   extendDBConfig,
	}
}

func initChatFileConfig() ChatFileConfig {

	fileServiceGrpcHost := os.Getenv("FILE_SERVICE_GRPC_HOST")
	fileServiceGrpcPort := os.Getenv("FILE_SERVICE_GRPC_PORT")
	messageServiceGrpcHost := os.Getenv("MESSAGE_SERVICE_GRPC_HOST")
	messageServiceGrpcPort := os.Getenv("MESSAGE_SERVICE_GRPC_PORT")

	org := os.Getenv("CHAT_FILE_ORG")
	flag := os.Getenv("CHAT_FILE_CRON_BATCH_FLAG")
	var batchFlag bool
	if flag == "true" {
		batchFlag = true
	} else {
		batchFlag = false
	}
	cron := os.Getenv("CHAT_FILE_CRON")

	batchConfig := &BatchConfig{
		Org:       org,
		BatchFlag: batchFlag,
		Cron:      cron,
	}

	return ChatFileConfig{

		FileServiceGrpcHost:    fileServiceGrpcHost,
		FileServiceGrpcPort:    fileServiceGrpcPort,
		MessageServiceGrpcHost: messageServiceGrpcHost,
		MessageServiceGrpcPort: messageServiceGrpcPort,
		ChatFileBatchConfig:    batchConfig,
	}

}

func ConnectDatabase(sfg *ServerConfig) (*gorm.DB, error) {

	//log.Println("env 읽음 " + sfg.dbConfig.Id + " : " + sfg.dbConfig.Pw + " : " + sfg.dbConfig.Host + " : " + sfg.dbConfig.Port + " : " + sfg.dbConfig.Database)
	dsn := sfg.dbConfig.Id + ":" + sfg.dbConfig.Pw + "@tcp(" + sfg.dbConfig.Host + ":" + sfg.dbConfig.Port + ")/" + sfg.dbConfig.Database + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}) // MYSQL

	if err != nil {
		return nil, err
	}

	//db.AutoMigrate(&sharedModels.ServiceUsers{})

	return db, nil
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

func NewFileServiceProtocolBufferClient(sfg *ServerConfig) (*grpc.ClientConn, error) {
	return grpc.NewClient(sfg.ChatFileConfig.FileServiceGrpcHost+":"+sfg.ChatFileConfig.FileServiceGrpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func NewMessageServiceProtocolBufferClient(sfg *ServerConfig) (*grpc.ClientConn, error) {
	return grpc.NewClient(sfg.ChatFileConfig.MessageServiceGrpcHost+":"+sfg.ChatFileConfig.MessageServiceGrpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
