package test

import (
	"multi-bot/constant"
	"multi-bot/utils/db_util"
	"multi-bot/utils/log"
	"os"
	"time"
)

func init() {
	time.Local = time.FixedZone("UTC", 0)
	initEnv()
	initConfig()

	//if controllerService := os.Getenv(constant.ControllerService); len(controllerService) > 0 {
	//	grpc_client.InitControllerClient(controllerService)
	//}
	//
	//if config.UseWebHook() {
	//	go grpc_server.Start(os.Getenv(constant.GrpcListenPort))
	//	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
	//		panic(err)
	//	}
	//} else {
	//	grpc_server.Start(os.Getenv(constant.GrpcListenPort))
	//}
}

func initEnv() {
	_ = os.Setenv("CONTROLLER_SERVICE=", "0.0.0.0:9092")
	_ = os.Setenv("DB_USER", "root")
	_ = os.Setenv("DB_PASSWORD", "root")
	_ = os.Setenv("DB_HOST", "10.11.95.56")
	_ = os.Setenv("DB_PORT", "3306")
	_ = os.Setenv("DB_NAME", "bot-my")
	_ = os.Setenv("ENV", "dev")
	_ = os.Setenv("LOG_CONSOLE_WRITE", "1")
	_ = os.Setenv("TSTORE_GRPC_SERVICE_ADDR", "0.0.0.0:2011")
}

func initConfig() {
	// connect to chain database
	dbHost := os.Getenv(constant.DB_HOST)
	dbPort := os.Getenv(constant.DB_PORT)
	dbName := os.Getenv(constant.DB_NAME)
	dbUser := os.Getenv(constant.DB_USER)
	dbPassword := os.Getenv(constant.DB_PASSWORD)
	err := db_util.ConnectDb(dbHost, dbPort, dbName, dbUser, dbPassword)
	if err != nil {
		log.Error().Fields(map[string]interface{}{
			"action": "connect to db",
			"error":  err.Error(),
		})
	}
}
