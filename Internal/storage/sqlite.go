package storage

import (
	"database/sql"

	"github.com/aadityya4real/Task-manager/internal/types"
)

type Store struct {
	DB *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{DB: db}
}

func (s *Store) InsertTask(t types.Task) (int64, error) {
	result, err := s.DB.Exec(
		"INSERT INTO tasks (title, done) VALUES (?, ?)",
		t.Title, false,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (s *Store) GetTasks() ([]types.Task, error) {
	rows, err := s.DB.Query("SELECT id, title, done FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []types.Task
	for rows.Next() {
		var t types.Task
		rows.Scan(&t.ID, &t.Title, &t.Done)
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (s *Store) UpdateTask(id int, t types.Task) error {
	_, err := s.DB.Exec(
		"UPDATE tasks SET title = ?, done = ? WHERE id = ?",
		t.Title, t.Done, id,
	)
	return err
}

func (s *Store) DeleteTask(id int) error {
	_, err := s.DB.Exec("DELETE FROM tasks WHERE id = ?", id)
	return err
}
