package runner

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"
)

type StepResult struct {
	Name     string
	Command  string
	Output   string
	Error    string
	Success  bool
	Duration time.Duration
}

func ExecuteSteps(steps []Step) ([]StepResult, error) {
	var results []StepResult


	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %w", err)
	}

	for _, step := range steps {
		fmt.Printf("▶ Running step: %s\n", step.Name)

		cmd := getShellCommand(step.Run)

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		start := time.Now()
		err := cmd.Run()
		duration := time.Since(start)

		success := err == nil

		logFileName := sanitizeFileName(step.Name) + ".log"
		logFilePath := filepath.Join("logs", logFileName)
		logData := fmt.Sprintf(
			"▶ STEP: %s\n⏱ Duration: %v\n\n🔹 STDOUT:\n%s\n🔸 STDERR:\n%s\n",
			step.Name, duration, stdout.String(), stderr.String(),
		)
		_ = os.WriteFile(logFilePath, []byte(logData), 0644)

		result := StepResult{
			Name:     step.Name,
			Command:  step.Run,
			Output:   stdout.String(),
			Error:    stderr.String(),
			Success:  success,
			Duration: duration,
		}

		if success {
			fmt.Printf("✅ Step '%s' completed in %v\n\n", step.Name, duration)
		} else {
			fmt.Printf("❌ Step '%s' failed in %v\n\n", step.Name, duration)
		}

		results = append(results, result)
	}

	return results, nil
}

func sanitizeFileName(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_-]+`)
	return re.ReplaceAllString(name, "_")
}
func getShellCommand(command string) *exec.Cmd {
	if isWindows() {
		return exec.Command("cmd", "/C", command)
	}
	return exec.Command("bash", "-c", command)
}

func isWindows() bool {
	return os.PathSeparator == '\\'
}
