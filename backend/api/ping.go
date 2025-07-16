package api

import (
	"net/http"
	"strings"

	"github.com/dhananjaykakade/mini-ci/backend/runner"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid ping URL", http.StatusBadRequest)
		return
	}
	containerID := parts[2]

	runner.UpdateContainerLastAccess(containerID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}
