package main

import (
	"Log-Injector/controllers/logController"
	"Log-Injector/db"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {

	// Create a new router using gorilla/mux
	r := mux.NewRouter()
	handler := func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "Hello from Server !!!")
		if err != nil {
			return
		}
	}

	// Initialize MongoDB connection
	err := db.Init()
	if err != nil {
		fmt.Println("Failed to connect to MongoDB:", err)
		return
	}

	// Initialize your MongoDB client
	mongoClient := db.GetClient()
	logsDB := db.NewMongoDB(mongoClient)

	// Initialize MongoDB collections
	err = logsDB.InitLogEntryCollection()
	if err != nil {
		fmt.Println("Failed to initialize logs collections:", err)
		return
	}

	// Start the Goroutine to process logs in the background
	go logController.ProcessLogs(logsDB, 1000)

	// routes
	r.HandleFunc("/", handler)

	r.HandleFunc("/api/logs", func(w http.ResponseWriter, r *http.Request) {
		logController.LogHandler(w, r)
	}).Methods("POST")

	r.HandleFunc("/api/search/filters", func(w http.ResponseWriter, r *http.Request) {
		logController.SearchLogsWithFilters(w, r, logsDB)
	}).Methods("GET")

	port := 3000
	fmt.Printf("Server is listening on :%d...\n", port)
	http.Handle("/", r)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}
