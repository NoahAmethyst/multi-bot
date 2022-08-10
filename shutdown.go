package main

import (
	"multi-bot/utils/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

//gracefulShutdown waits for termination syscalls and doing clean up operations after received it
func gracefulShutdown(server *http.Server) {
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGSTOP, syscall.SIGKILL, syscall.SIGHUP)
	go func() {
		sig := <-signalChannel
		log.Info().Msgf("shut down,sign:%v", sig)
		// code here for graceful shut down
	}()
}
