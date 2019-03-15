package main

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"taskmaster/id"
	"time"
)

type Workflow struct {
	Id string `json:"id"`
	Name string `json:"name"`
	WorkerMatchFunction string `json:"workerMatchFunction"`
	Stages []WorkflowStage `json:"stages"`
}

type WorkflowStage struct {
	Order int `json:"order"`
	WorkerMatchFunction string `json:"workerMatchFunction"`
	StageTimeout int64 `json:"stageTimeout"`
	SkipIfNoMatches bool `json:"skipIfNoMatches"`
}


type Message struct {
	TargetWorker string `json:"targetWorker"`
	Message map[string]interface{} `json:"message"`
}

const evaluationInterval = 5
const defaultStageTime = 180

func GetWorkflowsHandler(w http.ResponseWriter, r *http.Request) {
	workflows, _ := GetWorkflows()

	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.Encode(map[string]interface{}{
		"workflows": workflows,
	})
}

func GetWorkflows() ([]Workflow, error) {
	session := Store.Session.Clone()
	defer session.Close()

	var workflows []Workflow

	c := session.DB("taskmaster").C("workflows")
	err := c.Find(nil).All(&workflows)

	return workflows, err
}

func CreateWorkflow(workflow Workflow) (Workflow, error) {
	session := Store.Session.Clone()
	defer session.Close()

	for index, stage := range workflow.Stages {
		if stage.StageTimeout <= 0 {
			stage.StageTimeout = defaultStageTime
			workflow.Stages[index] = stage
		}
	}

	workflow.Id = id.GenerateId()

	c := session.DB("taskmaster").C("workflows")
	err := c.Insert(workflow)

	return workflow, err
}

func addTaskToWorkflow(workflow Workflow, task Task) {
	for _, stage := range workflow.Stages {
		start := time.Now()

		for time.Since(start).Seconds() < float64(stage.StageTimeout) {
			allWorkers, _ := GetWorkers()
			workflowWorkers, _ := MatchWorkers(allWorkers, workflow.WorkerMatchFunction, task)
			stageWorkers, _ := MatchWorkers(workflowWorkers, stage.WorkerMatchFunction, task)

			if len(stageWorkers) < 1 && stage.SkipIfNoMatches {
				logrus.Tracef("No workers found for Workflow %s Stage %d. Skipping stage.",
					workflow.Id, stage.Order)
				break
			}

			for _, worker := range stageWorkers {
				var reservationMessage = Message{
					TargetWorker: worker.Id,
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

		logrus.Tracef("Workflow %s stage %d timed out.", workflow.Id, stage.Order)
	}

	logrus.Tracef("Workflow %s timed out for task %s.", workflow.Id, task.Id)
}