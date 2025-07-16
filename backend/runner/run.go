package runner

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func RunDockerContainer(imageTag, containerName string, externalPort, internalPort int, env map[string]string) (string, error) {
	args := []string{"run", "-d"}

	args = append(args, "-p", fmt.Sprintf("%d:%d", externalPort, internalPort))

	for k, v := range env {
		args = append(args, "-e", fmt.Sprintf("%s=%s", k, v))
	}

	args = append(args, "--name", containerName, imageTag)

	cmd := exec.Command("docker", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("docker run failed: %v\n%s", err, string(out))
	}
	containerID := string(out)
	containerID = strings.TrimSpace(containerID)

	storeMutex.Lock()
	containerStore[containerID] = &ContainerInfo{
		ContainerID: containerID,
		Port:        externalPort,
		LastAccess:  time.Now(),
	}
	storeMutex.Unlock()

	return containerID, nil

}
