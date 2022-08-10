package config

import (
	"multi-bot/constant"
	"os"
)

func EnvIsDev() bool {
	return os.Getenv("ENV") == "dev"
}

func UseWebHook() bool {
	return len(os.Getenv(constant.USE_WEBHOOK)) > 0
}

func UseConsoleWrite() bool {
	return os.Getenv("LOG_CONSOLE_WRITE") == "1"
}

func GetPodName() string {
	return os.Getenv("POD_NAME")
}
