package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

type WorkerPayload struct {
	TaskID  string          `json:"task_id"`
	Slug    string          `json:"slug"`
	Payload json.RawMessage `json:"payload"`
}

type CsvPayload struct {
	SourceURI  string `json:"source_uri"`
	TargetURI  string `json:"target_uri"`
	Layout     string `json:"layout"`
	HasHeaders bool   `json:"has_headers"`
	FontSize   int    `json:"font_size"`
}

type WorkerResultMessage string

const (
	WorkerResultSuccessMesssage WorkerResultMessage = "success"
	WorkerResultFailedMessage   WorkerResultMessage = "failed"
	WorkerResultACKMessage      WorkerResultMessage = "ack"
)

type WorkerResult struct {
	TaskID        string              `json:"task_id"`
	ResultMessage WorkerResultMessage `json:"result_message"`
	Error         json.RawMessage     `json:"error,omitempty"`
	Timestamp     time.Time           `json:"timestamp,omitempty"`
}

var stdoutMu sync.Mutex

func writeResult(res WorkerResult) {
	stdoutMu.Lock()
	defer stdoutMu.Unlock()
	_ = json.NewEncoder(os.Stdout).Encode(res)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var wp WorkerPayload
		if err := json.Unmarshal(line, &wp); err != nil {
			log.Printf("Failed to decode worker payload wrapper: %v", err)
			continue
		}

		// 1. Immediately write ACK back to stdout
		writeResult(WorkerResult{
			TaskID:        wp.TaskID,
			ResultMessage: WorkerResultACKMessage,
			Timestamp:     time.Now(),
		})

		// 2. Process task asynchronously in a goroutine
		go func(taskID string, rawPayload json.RawMessage) {
			var payload CsvPayload
			if err := json.Unmarshal(rawPayload, &payload); err != nil {
				log.Printf("Failed to decode csv payload: %v", err)
				writeError(taskID, err.Error())
				return
			}

			// Simulate Work (e.g., generating PDF)
			time.Sleep(1 * time.Second) // Fake processing time

			writeResult(WorkerResult{
				TaskID:        taskID,
				ResultMessage: WorkerResultSuccessMesssage,
				Timestamp:     time.Now(),
			})
		}(wp.TaskID, wp.Payload)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading stdin: %v", err)
	}
}

func writeError(taskID string, errMsg string) {
	errBytes, _ := json.Marshal(errMsg)
	writeResult(WorkerResult{
		TaskID:        taskID,
		ResultMessage: WorkerResultFailedMessage,
		Error:         errBytes,
		Timestamp:     time.Now(),
	})
}
