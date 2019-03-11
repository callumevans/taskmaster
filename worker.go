package main

import (
	"encoding/json"
	"net/http"
	"taskmaster/id"
	"taskmaster/redis"
)

type Worker struct {
	Id string `json:"id"`
	Attributes interface{} `json:"attributes"`
}

func GetWorkersHandler(w http.ResponseWriter, r *http.Request) {
	workers, _ := GetWorkers()

	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)

	_ = encoder.Encode(map[string]interface{}{
		"workers": workers,
	})
}

func CreateWorkerHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var worker Worker
	err := decoder.Decode(&worker)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	created, err := CreateWorker(worker)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	encoder := json.NewEncoder(w)
	encoder.Encode(created)
}

func GetWorkers() ([]Worker, error) {
	res, err := redisConnection.JsonClient.JSONGet("workers", ".")

	if err != nil && err != redis.Nil {
		return nil, err
	}

	var workers []Worker
	_ = json.Unmarshal(res.([]byte), &workers)

	return workers, nil
}

func CreateWorker(worker Worker) (*Worker, error) {
	worker.Id = id.GenerateId()

	_, err := redisConnection.JsonClient.JSONArrAppend("workers", ".", worker)

	if err != nil {
		return nil, err
	}

	return &worker, nil
}