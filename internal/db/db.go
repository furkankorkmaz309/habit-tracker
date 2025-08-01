package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	//create and connect database
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		return nil, fmt.Errorf("DB_PATH env variable not set")
	}

	err := os.MkdirAll(dbPath, 0755)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while creating folder : %v", err)
	}

	db, err := sql.Open("sqlite3", dbPath+"data.db")
	if err != nil {
		return nil, fmt.Errorf("an error occurred while creating database : %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("an error occurred while connecting database : %v", err)
	}

	// create users table
	queryUsers := `
	CREATE TABLE IF NOT EXISTS users(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT NOT NULL UNIQUE,
	email TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL
	)`

	_, err = db.Exec(queryUsers)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while creating users table : %v", err)
	}

	// create habits table
	queryHabits := `
	CREATE TABLE IF NOT EXISTS habits(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
	description TEXT NOT NULL,
	frequency TEXT NOT NULL,
	day TEXT NOT NULL,
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
