package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dhananjaykakade/mini-ci/backend/runner"
)

func BuildStreamHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Setup headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Decode request body
	var config runner.DeployConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		fmt.Fprintf(w, "data: Error decoding request: %v\n\n", err)
		flusher.Flush()
		return
	}

	logChan := make(chan string)

	go func() {
	
		runner.ExecutePipelineWithLogs(config, logChan)
	}()

	for logLine := range logChan {
		fmt.Fprintf(w, "data: %s\n\n", logLine)
		flusher.Flush()
	}

	// Final success message
	fmt.Fprintf(w, "data: ✅ Deployment complete!\n\n")
	flusher.Flush()
}
func LogsHandler(w http.ResponseWriter, r *http.Request) {
	buildID := r.URL.Path[len("/logs/"):]
	logChan, ok := buildLogs[buildID]
	if !ok {
		http.Error(w, "Build not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	for msg := range logChan {
		fmt.Fprintf(w, "data: %s\n\n", msg)
		flusher.Flush()
	}

	fmt.Fprintf(w, "data: ✅ Logs complete for build %s\n\n", buildID)
	flusher.Flush() 
}
