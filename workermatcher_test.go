package main

import (
	"testing"
)

var workerTestData = []WorkerDto{
	{
		Id: "1",
		Attributes: map[string]interface{}{
			"skills": []string { "chat", "voice" },
		},
	},
	{
		Id: "2",
		Attributes: map[string]interface{}{
			"skills": []string { "voice" },
		},
	},
	{
		Id: "3",
		Attributes: map[string]interface{}{
			"skills": []string { "chat" },
		},
	},
	{
		Id: "4",
		Attributes: map[string]interface{}{
			"uniqueAttribute": "test",
		},
	},
}

func TestGivenWorkersCanFilterWithInterpreter(t *testing.T) {
	const filterScript = "worker.attributes.skills.indexOf('chat') > -1"

	result, err := MatchWorkers(workerTestData, filterScript, Task{})

	if err != nil {
		t.Fatalf("Did not expect an error but got %s", err.Error())
	}

	if len(result) != 2 {
		t.Fatalf("Expected 2 workers but got %d", len(result))
	}

	if result[0].Id != workerTestData[0].Id || result[1].Id != workerTestData[2].Id {
		t.Fatalf("Didn't match expected workers. \n Expected: %s \n Actual: %s", result, workerTestData)
	}
}

func TestGivenWorkersCanHandleMissingMatchCriteria(t *testing.T) {
	const filterScript = "worker.fakeproperty.skills.indexOf('chat') > -1"

	result, err := MatchWorkers(workerTestData, filterScript, Task{})

	if err != nil {
		t.Fatalf("Did not expect an error but got %s", err.Error())
	}

	if len(result) != 0 {
		t.Fatalf("Expected 0 workers but got %d", len(result))
	}
}

func TestGivenWorkersCanHandleSomeMissingMatchCriterias(t *testing.T) {
	const filterScript = "worker.attributes.uniqueAttribute == 'test'"

	result, err := MatchWorkers(workerTestData, filterScript, Task{})

	if err != nil {
		t.Fatalf("Did not expect an error but got %s", err.Error())
	}

	if len(result) != 1 {
		t.Fatalf("Expected 1 workers but got %d", len(result))
	}

	if result[0].Id != workerTestData[3].Id {
		t.Fatalf("Didn't match expected workers. \n Expected: %s \n Actual: %s", result, workerTestData)
	}
}

func TestMatchingFunctionIncludesTasks(t *testing.T) {
	const filterScript = "task.attributes.taskType == 'chat'"

	result, err := MatchWorkers(workerTestData, filterScript, Task{
		Attributes: map[string]interface{}{
			"taskType": "chat",
		},
	})

	if err != nil {
		t.Fatalf("Did not expect an error but got %s", err.Error())
	}

	if len(result) != len(workerTestData) {
		t.Fatalf("Expected 4 workers but got %d", len(result))
	}
}

func TestMatchingFunctionIncludesTasks_NoResults(t *testing.T) {
	const filterScript = "task.fakeProp == 'missing'"

	result, err := MatchWorkers(workerTestData, filterScript, Task{
		Attributes: map[string]interface{}{
			"taskType": "chat",
		},
	})

	if err != nil {
		t.Fatalf("Did not expect an error but got %s", err.Error())
	}

	if len(result) != 0 {
		t.Fatalf("Expected 0 workers but got %d", len(result))
	}
}