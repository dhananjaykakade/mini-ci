package runner

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

func BuildDockerImage(tag, contextPath string) error {
	cmd := exec.Command("docker", "build", "-t", tag, contextPath)
	cmd.Stdout = nil
	cmd.Stderr = nil
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker build failed: %v\n%s", err, string(out))
	}
	return nil
}

func ExecutePipelineWithLogs(config DeployConfig, logs chan<- string) {
	defer close(logs)

	id, path, err := CloneRepo(config.RepoURL)
	if err != nil {
		logs <- fmt.Sprintf("❌ Clone failed: %v", err)
		return
	}

	logs <- fmt.Sprintf("📁 Repo cloned to %s", path)
	defaults := GetDefaultsByAppType(config.AppType)

	dockerfileConfig := DockerfileConfig{
		Language:     defaults.Language,
		InstallCmd:   coalesce(config.InstallCmd, defaults.InstallCmd),
		BuildCmd:     coalesce(config.BuildCmd, defaults.BuildCmd),
		StartCmd:     coalesce(config.StartCmd, defaults.StartCmd),
		ExposePort:   ifZero(config.ExposePort, defaults.ExposePort),
		Env:          config.Env,
		OutputFolder: defaults.OutputFolder,
	}

	if err := GenerateDockerfile(path, dockerfileConfig); err != nil {
		logs <- fmt.Sprintf("❌ Dockerfile generation failed: %v", err)
		return
	}

	imageTag := "ci-" + id
	logs <- fmt.Sprintf("📦 Building Docker image: %s", imageTag)

	if err := BuildDockerImageWithLogs(imageTag, path, logs); err != nil {
		logs <- fmt.Sprintf("❌ Build failed: %v", err)
		return
	}

	port := GetRandomPort()
	containerID, err := RunDockerContainer(imageTag, "container-"+id, port, config.ExposePort, config.Env)
	if err != nil {
		logs <- fmt.Sprintf("❌ Run failed: %v", err)
		return
	}

	logs <- fmt.Sprintf("🚀 Deployed! Access at http://localhost:%d (container: %s)", port, containerID)
}

func BuildDockerImageWithLogs(tag, path string, logs chan<- string) error {
	cmd := exec.Command("docker", "build", "-t", tag, path)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return err
	}

	scan := func(r io.Reader) {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			logs <- scanner.Text()
		}
	}

	go scan(stdout)
	go scan(stderr)

	return cmd.Wait()
}
func coalesce(value, fallback string) string {
	if value != "" {
		return value
	}
	return fallback
}

func ifZero(val, fallback int) int {
	if val != 0 {
		return val
	}
	return fallback
}
