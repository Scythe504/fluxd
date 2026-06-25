package main

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type CsvPayload struct {
	SourceURI  string `json:"source_uri"`
	TargetURI  string `json:"target_uri"`
	Layout     string `json:"layout"`
	HasHeaders bool   `json:"has_headers"`
	FontSize   int    `json:"font_size"`
}

type WorkerResult struct {
	Status      string `json:"status"`
	Pages       int    `json:"pages"`
	TimeTakenMs int64  `json:"time_taken_ms"`
	Error       string `json:"error,omitempty"`
}

func main() {
	start := time.Now()

	// 1. Read payload from Stdin
	var payload CsvPayload
	if err := json.NewDecoder(os.Stdin).Decode(&payload); err != nil {
		log.Fatalf("Failed to decode stdin: %v", err)
	}

	// 2. Simulate Work (e.g., generating PDF)
	time.Sleep(1 * time.Second) // Fake processing time

	// 3. Write result to Stdout
	result := WorkerResult{
		Status:      "success",
		Pages:       14, // Simulated page count
		TimeTakenMs: time.Since(start).Milliseconds(),
	}

	if err := json.NewEncoder(os.Stdout).Encode(result); err != nil {
		log.Fatalf("Failed to write stdout: %v", err)
	}
}