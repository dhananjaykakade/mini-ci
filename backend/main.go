package main

import (
	"log"
	"net/http"

	"github.com/dhananjaykakade/mini-ci/backend/api"
	"github.com/dhananjaykakade/mini-ci/backend/runner"
)

func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		h(w, r)
	}
}

func main() {

	runner.StartCleanupWorker()


	http.HandleFunc("/test-stream", withCORS(api.TestStreamHandler))
	http.HandleFunc("/build-stream", withCORS(api.BuildStreamHandler))
	http.HandleFunc("/deploy", withCORS(api.DeployHandler))
	http.HandleFunc("/ping/", withCORS(api.PingHandler))
	http.HandleFunc("/logs/", withCORS(api.LogsHandler))
	http.HandleFunc("/health", withCORS(api.HealthHandler))

	log.Printf("Server started on :8080")
	http.ListenAndServe(":8080", nil)
	log.Printf("Server stopped")

}
