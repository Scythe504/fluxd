package main

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type VideoPayload struct {
	SourceURI    string `json:"source_uri"`
	TargetURI    string `json:"target_uri"`
	Resolution   string `json:"resolution"`
	Format       string `json:"format"`
	ExtractAudio bool   `json:"extract_audio"`
}

type WorkerResult struct {
	Status       string `json:"status"`
	BytesWritten int    `json:"bytes_written"`
	TimeTakenMs  int64  `json:"time_taken_ms"`
	Error        string `json:"error,omitempty"`
}

func main() {
	start := time.Now()

	// 1. Read payload from Stdin
	var payload VideoPayload
	if err := json.NewDecoder(os.Stdin).Decode(&payload); err != nil {
		log.Fatalf("Failed to decode stdin: %v", err)
	}

	// 2. Simulate Work (e.g., FFMPEG processing)
	time.Sleep(2 * time.Second) // Fake processing time

	// 3. Write result to Stdout
	result := WorkerResult{
		Status:       "success",
		BytesWritten: 10485760, // 10MB simulated
		TimeTakenMs:  time.Since(start).Milliseconds(),
	}

	if err := json.NewEncoder(os.Stdout).Encode(result); err != nil {
		log.Fatalf("Failed to write stdout: %v", err)
	}
}