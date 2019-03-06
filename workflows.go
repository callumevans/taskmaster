package main

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
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

const evaluationInterval = 5
const defaultStageTime = 180

func GetWorkflows() ([]Workflow, error) {
	res, err := GetJsonClient().JSONGet("workflows", ".")

	if err != nil && err != redis.Nil {
		return nil, err
	}

	var workflows []Workflow
	_ = json.Unmarshal(res.([]byte), &workflows)

	return workflows, nil
}

func CreateWorkflow(workflow Workflow) (*Workflow, error) {
	workflow.Id = GenerateId()

	for index, stage := range workflow.Stages {
		if stage.StageTimeout <= 0 {
			stage.StageTimeout = defaultStageTime
			workflow.Stages[index] = stage
		}
	}

	_, err := GetJsonClient().JSONArrAppend("workflows", ".", workflow)

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
			task.Id = GenerateId()
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
				logrus.Infof("No workers found for Workflow %s Stage %d. Skipping stage.",
					workflow.Id, stage.Order)
				break
			}

			for _, worker := range stageWorkers {
				logrus.Tracef("Worker %s pinged with task %s", worker.Id, task.Id)
			}

			time.Sleep(evaluationInterval * time.Second)
		}

		logrus.Tracef("Workflow %s stage %d timed out.", workflow.Id, stage.Order)
	}

	logrus.Tracef("Workflow %s timed out for task %s.", workflow.Id, task.Id)
}