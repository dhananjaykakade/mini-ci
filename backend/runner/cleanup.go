package runner

import (
	"log"
	"os/exec"
	"time"
)


func StartCleanupWorker() {
	go func() {
		for {
			time.Sleep(1 * time.Minute)

			now := time.Now()
			storeMutex.Lock()
			for id, info := range containerStore {
				if now.Sub(info.LastAccess) > 1*time.Minute {
					log.Printf("🧹 Stopping idle container: %s", id)
					stopAndRemoveContainer(id)
					delete(containerStore, id)
				}
			}
			storeMutex.Unlock()
		}
	}()
}

func stopAndRemoveContainer(containerID string) {
	exec.Command("docker", "stop", containerID).Run()
	exec.Command("docker", "rm", containerID).Run()
}
