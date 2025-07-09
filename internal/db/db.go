package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	//create and connect database
	err := os.MkdirAll("../../internal/data", 0755)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while creating folder : %v", err)
	}

	db, err := sql.Open("sqlite3", "../../internal/data/data.db")
	if err != nil {
		return nil, fmt.Errorf("an error occurred while creating database : %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("an error occurred while connecting database : %v", err)
	}

	// create users table

	// create habits table
	queryHabits := `
	CREATE TABLE IF NOT EXISTS habits(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
	description TEXT NOT NULL,
	frequency INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  	user_id INTEGER,
	FOREIGN KEY (user_id) REFERENCES users(id)
	)`

	_, err = db.Exec(queryHabits)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while creating habits table : %v", err)
	}

	// create habit_checkins table

	return db, nil
}
