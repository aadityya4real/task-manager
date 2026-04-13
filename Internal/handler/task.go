package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"context"

	"github.com/aadityya4real/Task-manager/internal/storage"
	"github.com/aadityya4real/Task-manager/internal/types"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func TaskHandler(store *storage.Store, rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// 🔹 POST → Create Task
		if r.Method == "POST" {
			var t types.Task

			err := json.NewDecoder(r.Body).Decode(&t)
			if err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}

			if t.Title == "" {
				http.Error(w, "Title is required", http.StatusBadRequest)
				return
			}

			id, err := store.InsertTask(t)
			if err != nil {
				http.Error(w, "Failed to insert", http.StatusInternalServerError)
				return
			}

			t.ID = int(id)
			t.Done = false

			json.NewEncoder(w).Encode(t)
			return
		}

		// 🔹 GET → Get all tasks
		if r.Method == "GET" {
			tasks, err := store.GetTasks()
			if err != nil {
				http.Error(w, "Failed to fetch", http.StatusInternalServerError)
				return
			}

			json.NewEncoder(w).Encode(tasks)
			return
		}

		// 🔹 PUT → Update Task
		if r.Method == "PUT" {
			idStr := r.URL.Query().Get("id")
			if idStr == "" {
				http.Error(w, "ID is required", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "Invalid ID", http.StatusBadRequest)
				return
			}

			var t types.Task
			err = json.NewDecoder(r.Body).Decode(&t)
			if err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}

			err = store.UpdateTask(id, t)
			if err != nil {
				http.Error(w, "Failed to update", http.StatusInternalServerError)
				return
			}

			json.NewEncoder(w).Encode(map[string]string{
				"status": "updated",
			})
			return
		}

		// 🔹 DELETE → Delete Task
		if r.Method == "DELETE" {
			idStr := r.URL.Query().Get("id")
			if idStr == "" {
				http.Error(w, "ID is required", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "Invalid ID", http.StatusBadRequest)
				return
			}

			err = store.DeleteTask(id)
			if err != nil {
				http.Error(w, "Failed to delete", http.StatusInternalServerError)
				return
			}

			json.NewEncoder(w).Encode(map[string]string{
				"status": "deleted",
			})
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
