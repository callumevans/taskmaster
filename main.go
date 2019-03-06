package main

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)

	InitialiseRedis()
	StartApi(5000)
}
