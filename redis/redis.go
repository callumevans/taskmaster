package redis

import (
	"github.com/go-redis/redis"
	"github.com/nitishm/go-rejson"
	"github.com/sirupsen/logrus"
)

type Connection struct {
	Client redis.Client
	JsonClient rejson.Handler
}

const Nil = redis.Nil

func InitialiseRedis() Connection {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})

	_, err := client.Ping().Result()

	if err != nil {
		logrus.Errorf("Could not connect to redis server: %s.", err.Error())
		panic(err)
	}

	redisJsonHandler := rejson.NewReJSONHandler()
	redisJsonHandler.SetGoRedisClient(client)

	initWorkers(*redisJsonHandler)
	initWorkflows(*redisJsonHandler)

	return Connection{
		Client: *client,
		JsonClient: *redisJsonHandler,
	}
}

func initWorkers(redisJsonClient rejson.Handler) {
	_, err := redisJsonClient.JSONGet("workers", ".")

	if err == redis.Nil {
		_, err = redisJsonClient.JSONSet("workers", ".", []string{})
	} else if err != nil {
		panic(err)
	}
}

func initWorkflows(redisJsonClient rejson.Handler) {
	_, err := redisJsonClient.JSONGet("workflows", ".")

	if err == redis.Nil {
		_, err = redisJsonClient.JSONSet("workflows", ".", []string{})
	} else if err != nil {
		panic(err)
	}
}