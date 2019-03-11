package main

import (
	"encoding/json"
	"github.com/go-redis/redis"
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

type Task struct {
	Id string `json:"id"`
	WorkflowId string `json:"workflowId"`
	Attributes interface{} `json:"attributes"`
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

func CreateWorkflowHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var workflow Workflow
	err := decoder.Decode(&workflow)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	created, err := CreateWorkflow(workflow)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	encoder := json.NewEncoder(w)
	encoder.Encode(created)
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var task Task
	err := decoder.Decode(&task)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	created, err := CreateTask(task)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if created == nil {
		http.Error(w, "Could not find a workflow to match task with", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	encoder := json.NewEncoder(w)
	encoder.Encode(created)
}

func GetWorkflows() ([]Workflow, error) {
	res, err := redisConnection.JsonClient.JSONGet("workflows", ".")

	if err != nil && err != redis.Nil {
		return nil, err
	}

	var workflows []Workflow
	_ = json.Unmarshal(res.([]byte), &workflows)

	return workflows, nil
}

func CreateWorkflow(workflow Workflow) (*Workflow, error) {
	workflow.Id = id.GenerateId()

	for index, stage := range workflow.Stages {
		if stage.StageTimeout <= 0 {
			stage.StageTimeout = defaultStageTime
			workflow.Stages[index] = stage
		}
	}

	_, err := redisConnection.JsonClient.JSONArrAppend("workflows", ".", workflow)

	if err != nil {
		return nil, err
	}

	return &workflow, nil
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
				redisConnection.Client.Publish("worker_reservations", string(messageJson))

				logrus.Tracef("Pinged worker %s with task %s", worker.Id, task.Id)
			}

			time.Sleep(evaluationInterval * time.Second)
		}

		logrus.Tracef("Workflow %s stage %d timed out.", workflow.Id, stage.Order)
	}

	logrus.Tracef("Workflow %s timed out for task %s.", workflow.Id, task.Id)
}