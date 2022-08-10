package main

import (
	"multi-bot/cluster/grpc_server"
	"multi-bot/config"
	"multi-bot/constant"
	"multi-bot/manager"
	"multi-bot/memory"
	"multi-bot/service"

	"multi-bot/utils/db_util"
	"multi-bot/utils/log"
	"net/http"
	"os"
	"time"
)

func main() {

	httpServer := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: nil,
	}
	time.Local = time.FixedZone("UTC", 0)

	gracefulShutdown(httpServer)

	initConfig()
	manager.StartAllBots()

	go func() {
		memory.InitMemory()
		service.InitJob()
	}()

	if config.UseWebHook() {
		go grpc_server.Start(os.Getenv(constant.GrpcListenPort))
		if err := httpServer.ListenAndServe(); err != nil {
			log.Warn().Msgf("http server listen error:%s", err.Error())
		}
	} else {
		grpc_server.Start(os.Getenv(constant.GrpcListenPort))
	}

	c := make(chan bool, 2)
	<-c

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
