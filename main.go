package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"os"
	"taskmaster/datastore"
	"taskmaster/pubsub"
	"taskmaster/websockets"
)

var Store *datastore.DataStore
var RedisClient *redis.Client

// Message Types
const (
	TaskReservationCreated 	= "task.reservation_created"
	TaskAccepted          	= "task.accepted"
	TaskWorkflowTimeout		= "task.workflow_timeout"
	TaskStageTimeout		= "task.stage_timeout"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.TraceLevel)

	Store = datastore.CreateStore()
	redisClient, err := pubsub.CreateRedis()

	if err != nil {
		logrus.Panicf("Error connecting to redis: %s", err)
	}

	RedisClient = redisClient

	hub := websockets.NewHub()

	hub.On(TaskAccepted, func(message websockets.InboundMessage) {
		fmt.Println("Task accepted")
	})

	go hub.Run()
	go websockets.ListenForMessages(RedisClient, hub)

	StartApi(5000, hub)
}