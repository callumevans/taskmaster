package main

import "taskmaster/id"

type Worker struct {
	Id string `json:"id"`
	Attributes map[string]interface{} `json:"attributes"`
}

func GetWorkers() ([]Worker, error) {
	session := Store.Session.Clone()
	defer session.Close()

	var workers []Worker

	c := session.DB("taskmaster").C("workers")
	err := c.Find(nil).All(&workers)

	return workers, err
}

func CreateWorker(worker Worker) (Worker, error) {
	session := Store.Session.Clone()
	defer session.Close()

	worker.Id = id.GenerateId()

	c := session.DB("taskmaster").C("workers")
	err := c.Insert(worker)

	return worker, err
}