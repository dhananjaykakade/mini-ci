package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dhananjaykakade/mini-ci/backend/runner"
	"github.com/google/uuid"
)

var (
	buildLogs = make(map[string]chan string)
)

type DeploymentRequest struct {
	RepoURL    string            `json:"repoUrl"`
	AppType    string            `json:"appType"`
	Env        map[string]string `json:"env"`
	InstallCmd string            `json:"installCmd,omitempty"`
	BuildCmd   string            `json:"buildCmd,omitempty"`
	StartCmd   string            `json:"startCmd,omitempty"`
	Port       int               `json:"port,omitempty"`
}

func DeployHandler(w http.ResponseWriter, r *http.Request) {
	var req DeploymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	buildID := uuid.NewString() 
	logChan := make(chan string, 100)
	buildLogs[buildID] = logChan
	url := ""

	go func() {
		defer close(logChan)
		defer delete(buildLogs, buildID)

		id, path, err := runner.CloneRepo(req.RepoURL)
		if err != nil {
			logChan <- "❌ Clone failed: " + err.Error()
			return
		}
		logChan <- "📁 Cloned to " + path

		
		def := runner.GetDefaultsByAppType(req.AppType)
		config := runner.DockerfileConfig{
			Language:     def.Language,
			InstallCmd:   fallback(req.InstallCmd, def.InstallCmd),
			BuildCmd:     fallback(req.BuildCmd, def.BuildCmd),
			StartCmd:     fallback(req.StartCmd, def.StartCmd),
			ExposePort:   fallbackInt(req.Port, def.ExposePort),
			OutputFolder: def.OutputFolder,
			Env:          req.Env,
		}

	
		if err := runner.GenerateDockerfile(path, config); err != nil {
			logChan <- "❌ Dockerfile generation failed: " + err.Error()
			return
		}
		logChan <- "📦 Dockerfile generated."

		imageTag := "ci-" + id
		logChan <- "🚀 Building image..."
		if err := runner.BuildDockerImageWithLogs(imageTag, path, logChan); err != nil {
			logChan <- "❌ Build failed: " + err.Error()
			return
		}

		port := runner.GetRandomPort()
		containerID, err := runner.RunDockerContainer(imageTag, "container-"+id, port, config.ExposePort, config.Env)
		if err != nil {
			logChan <- "❌ Run failed: " + err.Error()
			return
		}
		logChan <- fmt.Sprintf("✅ Running at http://localhost:%d (container: %s)", port, containerID)
		url = fmt.Sprintf("http://localhost:%d", port)
	}()

	json.NewEncoder(w).Encode(map[string]string{
		"buildID": buildID,
		"message": "Deployment started. successfully",
		"url":     url,
	})
}


func fallback(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func fallbackInt(a, b int) int {
	if a != 0 {
		return a
	}
	return b
}
