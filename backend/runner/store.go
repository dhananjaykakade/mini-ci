package runner

import (
	"sync"
	"time"
)

type ContainerInfo struct {
	ContainerID string
	Port        int
	LastAccess  time.Time
}

var containerStore = map[string]*ContainerInfo{}
var storeMutex = sync.Mutex{}

func UpdateContainerLastAccess(containerID string) {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	if info, ok := containerStore[containerID]; ok {
		info.LastAccess = time.Now()
	}
}
