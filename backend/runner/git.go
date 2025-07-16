package runner

import (
	"fmt"
	"os/exec"
)

func CloneRepo(repoURL string) (string, string, error) {
	id, path, err := GenerateWorkspace()
	if err != nil {
		return "", "", fmt.Errorf("failed to create workspace: %w", err)
	}

	cmd := exec.Command("git", "clone", repoURL, path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", "", fmt.Errorf("git clone failed: %v\n%s", err, string(output))
	}

	return id, path, nil
}
