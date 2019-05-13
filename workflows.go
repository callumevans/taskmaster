package main

import (
	"encoding/json"
	"net/http"
	"taskmaster/id"
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