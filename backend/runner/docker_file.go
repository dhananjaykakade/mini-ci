package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GenerateDockerfile(path string, config DockerfileConfig) error {
	lines := []string{}

	
	isNodeApp := config.Language == "node"


	tsconfigPath := filepath.Join(path, "tsconfig.json")
	tsExists := fileExists(tsconfigPath)

	
	if isNodeApp && !tsExists && isDefaultBuild(config.AppType, config.BuildCmd) {
		config.BuildCmd = ""
	}

	switch config.Language {
	case "node":
		lines = append(lines,
			"FROM node:22-alpine",
			"WORKDIR /app",
			"COPY . .",
		)

		if config.InstallCmd != "" {
			lines = append(lines, fmt.Sprintf("RUN %s", config.InstallCmd))
		}
		if config.BuildCmd != "" {
			lines = append(lines, fmt.Sprintf("RUN %s", config.BuildCmd))
		}

		for k, v := range config.Env {
			lines = append(lines, fmt.Sprintf("ENV %s=%s", k, v))
		}

		lines = append(lines,
			fmt.Sprintf("EXPOSE %d", config.ExposePort),
			fmt.Sprintf("CMD [\"sh\", \"-c\", \"%s\"]", config.StartCmd),
		)

	case "python":
		lines = append(lines,
			"FROM python:3.9",
			"WORKDIR /app",
			"COPY . .",
			fmt.Sprintf("RUN %s", config.InstallCmd),
			fmt.Sprintf("EXPOSE %d", config.ExposePort),
			fmt.Sprintf("CMD [\"sh\", \"-c\", \"%s\"]", config.StartCmd),
		)

	case "go":
		lines = append(lines,
			"FROM golang:1.20",
			"WORKDIR /app",
			"COPY . .",
		)

		if config.BuildCmd != "" {
			lines = append(lines, fmt.Sprintf("RUN %s", config.BuildCmd))
		}

		lines = append(lines,
			fmt.Sprintf("EXPOSE %d", config.ExposePort),
			fmt.Sprintf("CMD [\"sh\", \"-c\", \"%s\"]", config.StartCmd),
		)
	}

	dockerfilePath := filepath.Join(path, "Dockerfile")
	content := strings.Join(lines, "\n") + "\n"
	return os.WriteFile(dockerfilePath, []byte(content), 0644)
}


func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}


func isDefaultBuild(AppType, buildCmd string) bool {
	defaults := map[string]string{
		"react":  "npm run build",
		"nextjs": "npm run build",
		"vite":  "npm run build",
		"node":   "",
	}
	return buildCmd == defaults[AppType]
}
