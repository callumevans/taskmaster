package main

import (
	"encoding/json"
	"fmt"
	"github.com/robertkrimen/otto"
	"github.com/sirupsen/logrus"
)

func MatchWorkers(workers []Worker, matchFunction string, task Task) ([]Worker, error) {
	workersJson, _ := json.Marshal(workers)
	taskJson, _ := json.Marshal(task)

	return matchWorkersFromJson(string(workersJson), string(taskJson), matchFunction)
}

func matchWorkersFromJson(workerJson string, taskJson string, matchFunction string) ([]Worker, error) {
	javascript := fmt.Sprintf(`
		var workers = JSON.parse('%s');
		var task = JSON.parse('%s');
		var matches = [];

		for (var i = 0; i < workers.length; i++) {
			var worker = workers[i];

			try {
				if (%s) {
					matches.push(worker);
				}
			} finally {
				continue;
			}
		}
		
		matches = JSON.stringify(matches);`,
	workerJson, taskJson, matchFunction)

	vm := otto.New()

	_, err := vm.Run(javascript)

	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}

	val, _ := vm.Get("matches")
	strVal, _ := val.ToString()

	var workers []Worker
	json.Unmarshal([]byte(strVal), &workers)

	return workers, nil
}
