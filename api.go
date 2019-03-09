package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"taskmaster/websockets"
	"time"
)

func StartApi(port int, hub *websockets.Hub) {
	r := mux.NewRouter()

	r.HandleFunc("/workers", GetWorkersHandler).Methods(http.MethodGet)
	r.HandleFunc("/workers", CreateWorkerHandler).Methods(http.MethodPost)

	r.HandleFunc("/workflows", GetWorkflowsHandler).Methods(http.MethodGet)
	r.HandleFunc("/workflows", CreateWorkflowHandler).Methods(http.MethodPost)

	r.HandleFunc("/tasks", CreateTaskHandler).Methods(http.MethodPost)

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websockets.ServeWs(hub, w, r)
	})

	r.HandleFunc("/healthz", HealthzHandler)

	log.Info("Starting API")

	server := &http.Server{
		Addr: ":" + strconv.Itoa(port),
		Handler: logging()(r),
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 15 * time.Second,
	}

	log.Infof("API listening on port %d", port)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("API failed to start: %s", err)
	}
}

func GetWorkersHandler(w http.ResponseWriter, r *http.Request) {
	workers, _ := GetWorkers()

	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.Encode(map[string]interface{}{
		"workers": workers,
	})
}

func CreateWorkerHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var worker Worker
	err := decoder.Decode(&worker)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	created, err := CreateWorker(worker)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	encoder := json.NewEncoder(w)
	encoder.Encode(created)
}

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

func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func logging() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				log.Infof("%s %s %s %s", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()

			next.ServeHTTP(w, r)
		})
	}
}