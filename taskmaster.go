package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"taskmaster/api"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	api.Listen(5000)
}
