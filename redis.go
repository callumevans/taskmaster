package main

import (
	"github.com/go-redis/redis"
	"github.com/nitishm/go-rejson"
	"github.com/sirupsen/logrus"
)

var client *redis.Client
var redisJsonHandler *rejson.Handler

const Nil = redis.Nil

func InitialiseRedis() {
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})

	_, err := client.Ping().Result()

	if err != nil {
		logrus.Errorf("Could not connect to redis server: %s.", err.Error())
		panic(err)
	}

	redisJsonHandler = rejson.NewReJSONHandler()
	redisJsonHandler.SetGoRedisClient(client)

	initWorkers()
	initWorkflows()
}

func GetClient() *redis.Client {
	return client
}

func GetJsonClient() *rejson.Handler {
	return redisJsonHandler
}

func initWorkers() {
	_, err := redisJsonHandler.JSONGet("workers", ".")

	if err == Nil {
		_, err = redisJsonHandler.JSONSet("workers", ".", []string{})
	} else if err != nil {
		panic(err)
	}
}

func initWorkflows() {
	_, err := redisJsonHandler.JSONGet("workflows", ".")

	if err == Nil {
		_, err = redisJsonHandler.JSONSet("workflows", ".", []string{})
	} else if err != nil {
		panic(err)
	}
}