

package api

import (
	"net/http"

	"github.com/dhananjaykakade/mini-ci/backend/runner"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if !runner.CheckDockerRunning() {
		http.Error(w, "Docker is not running", 500)
		return
	}
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
