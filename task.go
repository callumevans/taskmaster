package main

type Task struct {
	Id string `json:"id"`
	WorkflowId string `json:"workflowId"`
	Attributes interface{} `json:"attributes"`
}