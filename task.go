package main

import "taskmaster/id"

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