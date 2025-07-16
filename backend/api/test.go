package api

import (
	"fmt"
	"net/http"
	"time"
)

func TestStreamHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	
	for i := 1; i <= 10; i++ {
		fmt.Fprintf(w, "data: Log line %d\n\n", i)
		flusher.Flush()
		time.Sleep(1 * time.Second)
	}

	fmt.Fprintf(w, "data: ✅ Done streaming test logs\n\n")
	flusher.Flush()
}
