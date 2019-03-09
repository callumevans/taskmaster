package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"taskmaster/websockets"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)

	InitialiseRedis()

	hub := websockets.NewHub()
	go hub.Run()
	go websockets.ListenForReservations(client, hub)

	StartApi(5000, hub)
}
