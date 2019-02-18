package redis

import "github.com/go-redis/redis"

var client *redis.Client

func init() {
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})

	_, err := client.Ping().Result()

	if err != nil {
		panic(err)
	}
}

func GetClient() *redis.Client {
	return client
}

func GetAllKeysInNamespace(namespace string) []string {
	var vals []string
	iterator := client.HScan(namespace, 0, "", 10).Iterator()

	for iterator.Next() {
		vals = append(vals, iterator.Val())

		if iterator.Err() != nil {
			panic(iterator.Err())
		}
	}

	return vals
}