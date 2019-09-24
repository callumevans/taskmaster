package main

import (
	"encoding/json"
	"github.com/globalsign/mgo/bson"
	"github.com/sirupsen/logrus"
	"taskmaster/id"
	"taskmaster/websockets"
	"time"
)

type Task struct {
	Id string `json:"id"`
	WorkflowId string `json:"workflowId"`
	Attributes interface{} `json:"attributes"`
}

func CreateTask(task Task) (*Task, error) {
	workflows, err := GetWorkflows()

	if err != nil {
		return nil, err
	}

	for _, workflow := range workflows {
		if workflow.Id == task.WorkflowId {
			task.Id = id.GenerateId()
			go addTaskToWorkflow(workflow, task)
			return &task, nil
		}
	}

	return nil, nil
}

func AddTaskToWorker(workerId string, task Task) error {
	session := Store.Session.Clone()
	defer session.Close()

	c := session.DB("taskmaster").C("workers")

	query := bson.M{"id": workerId}
	err := c.Update(query, bson.M{"$push": bson.M{"tasks": task}})

	return err
}

func RemoveTaskFromWorker(workerId string, taskId string) error {
	session := Store.Session.Clone()
	defer session.Close()

	c := session.DB("taskmaster").C("workers")

	query := bson.M{"id": workerId}
	err := c.Update(query, bson.M{"$pull": bson.M{"tasks": bson.M{"id": taskId}}})

	return err
}

func addTaskToWorkflow(workflow Workflow, task Task) {
	for _, stage := range workflow.Stages {
		start := time.Now()

		for time.Since(start).Seconds() < float64(stage.StageTimeout) {
			allWorkers, _ := GetWorkers()
			workflowWorkers, _ := MatchWorkers(allWorkers, workflow.WorkerMatchFunction, task)
			stageWorkers, _ := MatchWorkers(workflowWorkers, stage.WorkerMatchFunction, task)

			if len(stageWorkers) < 1 {
				logrus.Tracef("No workers found for Workflow %s Stage %d. Task: %s",
					workflow.Id, stage.Order, task.Id)

				if stage.SkipIfNoMatches {
					logrus.Tracef("Workflow %s skipping Stage %d. Task: %s",
						workflow.Id, stage.Order, task.Id)

					break
				}
			}

			for _, worker := range stageWorkers {
				var reservationMessage = websockets.OutboundMessage{
					TargetWorker: worker.Id,
					MessageType: TaskReservationCreated,
					Message: map[string]interface{}{
						"Task": task,
					},
				}

				var messageJson, _ = json.Marshal(reservationMessage)
				RedisClient.Publish("worker_reservations", string(messageJson))
				logrus.Tracef("Matched worker %s with task %s", worker.Id, task.Id)
			}

			time.Sleep(evaluationInterval * time.Second)
		}

		logrus.Tracef("Workflow %s stage %d timed out for task %s.", workflow.Id, stage.Order, task.Id)
	}

	logrus.Tracef("Workflow %s timed out for task %s.", workflow.Id, task.Id)
}