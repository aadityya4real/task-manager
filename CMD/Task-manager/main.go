package main

import (
	"encoding/json"
	"net/http"
)

// Task structure
type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// temporary storage (in-memory)
var tasks []Task
var currentID = 1

func taskHandler(w http.ResponseWriter, r *http.Request) {

	// POST → Create Task
	if r.Method == "POST" {
		var t Task

		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if t.Title == "" {
			http.Error(w, "Title is required", http.StatusBadRequest)
			return
		}

		t.ID = currentID
		currentID++
		tasks = append(tasks, t)

		json.NewEncoder(w).Encode(t)
		return
	}

	// GET → Get all tasks
	if r.Method == "GET" {
		json.NewEncoder(w).Encode(tasks)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func main() {

	http.HandleFunc("/task", taskHandler)

	http.ListenAndServe(":8080", nil)
}
