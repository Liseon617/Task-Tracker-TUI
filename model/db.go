package model

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB() error {
	var err error
	db, err = sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		status INTEGER,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`)
	if err != nil {
		return fmt.Errorf("Failed to create table: %v", err)
	}
	return nil
}

func SaveTask(t Task) error {
    tx, err := db.Begin()
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback() // Safe to call if tx.Commit() succeeds

    _, err = tx.Exec(
        "INSERT INTO tasks (title, description, status) VALUES (?, ?, ?)",
        t.title, t.description, t.status,
    )
    if err != nil {
        return fmt.Errorf("failed to insert task: %w", err)
    }

    return tx.Commit()
}

func LoadTasks() ([]Task, error) {
	rows, err := db.Query("SELECT title, description, status FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var title, description string
		var status Status
		if err := rows.Scan(&title, &description, &status); err != nil {
			return nil, err
		}
		tasks = append(tasks, Task{
			title: title, 
			description: description, 
			status: status,
		})
	}
	return tasks, nil
}

func UpdateTask (oldTask Task, newTitle, newDescription string) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	_, err = tx.Exec(
        `UPDATE tasks 
         SET title = ?, description = ?
         WHERE title = ? AND description = ?`,
        newTitle, newDescription,
        oldTask.title, oldTask.description, oldTask.status,
    )

	if err != nil {
		return fmt.Errorf("failed to update task %w", err)
	}

	return tx.Commit()
}

func DeleteTask(t Task) error {
    // First try to delete by ID if you add that field
    // Otherwise use a transaction to ensure atomic operation
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback() // Will be ignored if tx.Commit() succeeds
    
    _, err = tx.Exec("DELETE FROM tasks WHERE title = ? AND description = ? AND status = ?", 
        t.title, t.description, t.status)
    if err != nil {
        return err
    }
    
    return tx.Commit()
}
func ClearAllTasks() error {
    tx, err := db.Begin()
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback()

    _, err = tx.Exec("DELETE FROM tasks")
    if err != nil {
        return fmt.Errorf("failed to clear tasks: %w", err)
    }

    return tx.Commit()
}

