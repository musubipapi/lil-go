package dbUtils

import (
	"database/sql"
	"fmt"
	. "lil-todo/pkg/types"
	"log"
)

func ConnectToDb(filename ...string) (*sql.DB, error) {
	defaultFileName := "todo.db"
	dbFileName := defaultFileName

	if len(filename) > 0 && filename[0] != "" {
		dbFileName = filename[0]
	}

	db, err := sql.Open("sqlite", dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	// Check if the connection is successful
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// Create todo table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			task TEXT NOT NULL,
			completed BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func AddTodo(db *sql.DB, todo string) error {
	_, err := db.Exec(`
		INSERT INTO todos (task) VALUES (?)
	`, todo)

	return err
}

func DeleteTodo(db *sql.DB, id int) error {
	_, err := db.Exec(`
		DELETE FROM todos WHERE id = ?
	`, id)

	return err
}

func ToggleTodo(db *sql.DB, id int) (Todo, error) {
	var todo Todo
	err := db.QueryRow(`
	UPDATE todos
	SET completed = CASE WHEN completed = TRUE THEN FALSE ELSE TRUE END
	WHERE id = ? 
	RETURNING id, task, completed, created_at
	`, id).Scan(&todo.ID, &todo.Task, &todo.Completed, &todo.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return Todo{}, fmt.Errorf("no todo found with id %d", id)
		}
	}
	return todo, err
}

func ListAllTodos(db *sql.DB) ([]Todo, error) {
	rows, err := db.Query(`
		SELECT id, task, completed, created_at FROM todos
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Task, &todo.Completed, &todo.CreatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return todos, nil
}
