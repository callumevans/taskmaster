package main

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"taskmaster/id"
)

type Worker struct {
	Id string `json:"id"`
	Attributes interface{} `json:"attributes"`
}

func GetWorkers() ([]Worker, error) {
	res, err := GetJsonClient().JSONGet("workers", ".")

	if err != nil && err != redis.Nil {
		return nil, err
	}

	var workers []Worker
	_ = json.Unmarshal(res.([]byte), &workers)

	return workers, nil
}

func CreateWorker(worker Worker) (*Worker, error) {
	worker.Id = id.GenerateId()

	_, err := GetJsonClient().JSONArrAppend("workers", ".", worker)

	if err != nil {
		return nil, err
	}

	return &worker, nil
}