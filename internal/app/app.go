package app

import (
	"database/sql"
	"log"
)

type App struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}
