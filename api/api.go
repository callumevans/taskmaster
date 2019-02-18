package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"taskmaster/models"
	"time"
)


func Listen(port int) {
	r := mux.NewRouter()

	r.HandleFunc("/workers", GetWorkers).Methods(http.MethodGet)
	r.HandleFunc("/workers", CreateWorker).Methods(http.MethodPost)

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

func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func GetWorkers(w http.ResponseWriter, r *http.Request) {
	workers := models.GetWorkers()

	encoder := json.NewEncoder(w)
	encoder.Encode(workers)
}

func CreateWorker(w http.ResponseWriter, r *http.Request) {
	models.CreateNewWorker(map[string]interface{}{
		"tasks": "holla holla",
	})
	w.WriteHeader(http.StatusCreated)
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