package runner

import (
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

func GenerateWorkspace() (string, string, error) {
	id := uuid.New().String()[:6]
	if _, err := os.Stat("ci"); os.IsNotExist(err) {
		if err := os.Mkdir("ci", 0755); err != nil {
			return "", "", err
		}
	}
	path := filepath.Join("ci", "workspace-"+id)
	err := os.MkdirAll(path, 0755)
	return id, path, err
}

func GetRandomPort() int {
	rand.Seed(time.Now().UnixNano())
	return 8000 + rand.Intn(1000) 
}

func CheckDockerRunning() bool {
	cmd := exec.Command("docker", "info")
	err := cmd.Run()
	return err == nil
}

func DetectIfTSProject(path string) bool {
	tsconfigPath := filepath.Join(path, "tsconfig.json")
	_, err := os.Stat(tsconfigPath)
	return err == nil
}
