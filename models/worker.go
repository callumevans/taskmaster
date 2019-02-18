package models

import (
	_redis "github.com/go-redis/redis"
	"taskmaster/id"
	"taskmaster/redis"
)

const keyspace = "workers"

func CreateNewWorker(fields map[string]interface{}) {
	redis.GetClient().HMSet(keyspace + ":" + id.GenerateId(), map[string]interface{}{
		"tasks": "lololol",
	})
}

func GetWorkers() []string {
	pipe := redis.GetClient().Pipeline()
	defer pipe.Close()

	pipe.Scan(0, keyspace + ":*", 0)

	cmds, err := pipe.Exec()

	if err != nil {
		panic(err)
	}

	results, _, _ := cmds[0].(*_redis.ScanCmd).Result()

	return results
}

func GetWorker(id string) map[string]string {
	result := redis.GetClient().HGetAll(id)

	if result.Err() != nil {
		panic(result.Err())
	}

	return result.Val()
}