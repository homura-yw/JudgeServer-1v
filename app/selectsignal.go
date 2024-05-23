package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func register(UUID string) {
	redisClient := newClient()
	redisClient.Set(UUID, 0, 0)

	sigChan := make(chan os.Signal, 10)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-sigChan

	log.Println("service shutdown:", sig)
	redisClient.Del(UUID)
	os.Exit(1)
}
