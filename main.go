package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"taskmaster/redis"
	"taskmaster/websockets"
)

var redisConnection redis.Connection

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)

	redisConnection = redis.InitialiseRedis()

	hub := websockets.NewHub()

	go hub.Run()
	go websockets.ListenForReservations(redisConnection, hub)

	StartApi(5000, hub)
}