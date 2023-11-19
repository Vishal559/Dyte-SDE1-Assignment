package logController

import (
	"Log-Injector/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// LogsDatabase interface for handling database operations
type LogsDatabase interface {
	InsertLogs(logs []models.LogEntry) error
	SearchLogsWithFilters(query string, filters map[string]string, page int, pageSize int) (string, error)
}

// LogEntry struct represents the structure of the log message
type LogEntry struct {
	Level      string    `json:"level"`
	Message    string    `json:"message"`
	ResourceID string    `json:"resourceId"`
	Timestamp  time.Time `json:"timestamp"`
	TraceID    string    `json:"traceId"`
	SpanID     string    `json:"spanId"`
	Commit     string    `json:"commit"`
	Metadata   Metadata  `json:"metadata"`
}

// Metadata struct represents the metadata section of the log message
type Metadata struct {
	ParentResourceID string `json:"parentResourceId"`
}

var (
	logQueue = make(chan LogEntry, 100000000) // Buffered channel to store logs with a capacity of 100
	wg       sync.WaitGroup
)

// LogHandler function handles incoming POST requests to the /logs endpoint
func LogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var logMessage LogEntry
	err := json.NewDecoder(r.Body).Decode(&logMessage)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Send the log to the channel
	logQueue <- logMessage

	// Send a response back to the client
	w.WriteHeader(http.StatusAccepted)
	_, err = w.Write([]byte("Log request accepted successfully\n"))
	if err != nil {
		return
	}
}
func SearchLogsWithFilters(w http.ResponseWriter, r *http.Request, logsDB LogsDatabase) {
	query := r.URL.Query().Get("query")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	// Define the number of entries per page
	pageSize := 50

	// Extract filters from the request
	filters := make(map[string]string)
	if r.URL.Query().Get("level") != "" {
		filters["level"] = r.URL.Query().Get("level")
	}
	if r.URL.Query().Get("message") != "" {
		filters["message"] = r.URL.Query().Get("message")
	}
	if r.URL.Query().Get("resourceId") != "" {
		filters["resourceId"] = r.URL.Query().Get("resourceId")
	}
	if r.URL.Query().Get("timestamp") != "" {
		filters["timestamp"] = r.URL.Query().Get("timestamp")
	}
	if r.URL.Query().Get("traceId") != "" {
		filters["traceId"] = r.URL.Query().Get("traceId")
	}
	if r.URL.Query().Get("spanId") != "" {
		filters["spanId"] = r.URL.Query().Get("spanId")
	}
	if r.URL.Query().Get("commit") != "" {
		filters["commit"] = r.URL.Query().Get("commit")
	}
	if r.URL.Query().Get("parentResourceId") != "" {
		filters["parentResourceId"] = r.URL.Query().Get("parentResourceId")
	}

	// Perform text-based search with filters and pagination
	results, err := logsDB.SearchLogsWithFilters(query, filters, page, pageSize)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = fmt.Fprintf(w, `{"results": %v}`, results)
	if err != nil {
		fmt.Println("Error writing response:", err)
	}
}

// ProcessLogs function continuously processes logs from the channel
func ProcessLogs(logsDB LogsDatabase, batchSize int) {

	var batch []LogEntry
	// Create a channel to receive results from Goroutines
	resultChan := make(chan error)

	for {
		select {
		case logEntry, ok := <-logQueue:
			if !ok {
				// logQueue is closed, process any remaining logs in the current batch.
				if len(batch) > 0 {
					// Spawn a Goroutine to process the batch.
					wg.Add(1)
					go func(batch []LogEntry) {
						defer wg.Done()
						processLogBatch(logsDB, batch, resultChan)
					}(batch)
				}
				// Exit the function when logQueue is closed.
				return
			}

			// Accumulate log entries until reaching the batch size.
			batch = append(batch, logEntry)

			if len(batch) >= batchSize {
				// Spawn a Goroutine to process the batch.
				wg.Add(1)
				go func(batch []LogEntry) {
					defer wg.Done()
					processLogBatch(logsDB, batch, resultChan)
				}(batch)
				// Reset the batch.
				batch = nil
			}
		case <-time.After(2 * time.Second):
			// Process the batch even if it hasn't reached the batch size but a certain time has passed.
			if len(batch) > 0 {
				// Spawn a Goroutine to process the remaining batch.
				wg.Add(1)
				go func(batch []LogEntry) {
					defer wg.Done()
					processLogBatch(logsDB, batch, resultChan)
				}(batch)
				// Reset the batch.
				batch = nil
			}
		case err := <-resultChan:
			if err != nil {
				fmt.Printf("Error processing log batch: %s\n", err)
			}
		}

	}
}

func processLogBatch(logsDB LogsDatabase, batch []LogEntry, resultChan chan<- error) {
	// Perform processing on the log batch
	var dbBatchLogs []models.LogEntry

	for _, log := range batch {
		dbLog := models.LogEntry{
			Level:      log.Level,
			Message:    log.Message,
			ResourceId: log.ResourceID,
			Timestamp:  time.Time{},
			TraceId:    log.TraceID,
			SpanId:     log.SpanID,
			Commit:     log.Commit,
			Metadata: models.Metadata{
				ParentResourceId: log.Metadata.ParentResourceID,
			},
		}
		dbBatchLogs = append(dbBatchLogs, dbLog)
	}

	err := logsDB.InsertLogs(dbBatchLogs)
	if err != nil {
		resultChan <- err
	}

	resultChan <- nil
}
