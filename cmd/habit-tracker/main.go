package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/furkankorkmaz309/habit-tracker/internal/app"
	"github.com/furkankorkmaz309/habit-tracker/internal/db"
	"github.com/furkankorkmaz309/habit-tracker/internal/routes"
	"github.com/joho/godotenv"
)

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ltime|log.Ldate)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ltime|log.Ldate|log.Lshortfile)

	port := flag.String("http port", ":8080", "new http port")
	flag.Parse()

	db, err := db.InitDB()
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	err = godotenv.Load("../../.env")
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &app.App{
		DB:       db,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
	}

	r := routes.Routes(app)

	app.InfoLog.Printf("Server running on port %v\n", *port)
	http.ListenAndServe(*port, r)
}
