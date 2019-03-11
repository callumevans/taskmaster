package main

import (
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

	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Infof("Starting API on port %d", port)

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