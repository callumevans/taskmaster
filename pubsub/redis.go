package pubsub

import "github.com/go-redis/redis"

func CreateRedis() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		Password: "",
		DB: 0,
	})

	_, err := client.Ping().Result()

	return client, err
}
