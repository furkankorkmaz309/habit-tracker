package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/furkankorkmaz309/habit-tracker/internal/app"
	"github.com/furkankorkmaz309/habit-tracker/internal/models"
	"github.com/go-chi/chi"
)

func AddHabit(app app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//take habit input
		var habit models.Habit
		err := json.NewDecoder(r.Body).Decode(&habit)
		if err != nil {
			app.ErrorLog.Printf("an error occurred while decoding input : %v", err)
			http.Error(w, "Decode error", http.StatusInternalServerError)
			return
		}

		// check frequency
		fq, err := FrequencyConvert(habit.Frequency)
		if err != nil {
			app.ErrorLog.Printf("an error occurred while converting frequency : %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// add habit to database
		query := `INSERT INTO habits (title, description, frequency) VALUES (?, ?, ?)`
		result, err := app.DB.Exec(query, habit.Title, habit.Description, fq)
		if err != nil {
			app.ErrorLog.Printf("an error occurred while inserting habit to database : %v", err)
			http.Error(w, "DB insert error", http.StatusInternalServerError)
			return
		}

		id64, err := result.LastInsertId()
		if err != nil {
			app.ErrorLog.Printf("an error occurred while retrieve ID : %v", err)
			http.Error(w, "ID retrieve error", http.StatusInternalServerError)
			return
		}
		habit.ID = int(id64)

		// response habit and status code
		err = responseSuccess(w, habit, "Habit added successfully!", http.StatusCreated)
		if err != nil {
			app.ErrorLog.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GetHabits(app app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// create slice
		var habits []models.Habit

		// take values from db
		query := `SELECT id, title, description, frequency, created_at FROM habits`
		rows, err := app.DB.Query(query)
		if err != nil {
			app.ErrorLog.Printf("an error occurred while taking habits from database : %v", err)
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

		// add to slice
		for rows.Next() {
			var habit models.Habit
			err = rows.Scan(&habit.ID, &habit.Title, &habit.Description, &habit.Frequency, &habit.CreatedAt)
			if err != nil {
				app.ErrorLog.Printf("an error occurred while scanning row : %v", err)
				http.Error(w, "Row scan error", http.StatusInternalServerError)
				return
			}

			habits = append(habits, habit)
		}

		// response
		err = responseSuccess(w, habits, "Habits listed successfully!", http.StatusOK)
		if err != nil {
			app.ErrorLog.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GetHabit(app app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// take id from url
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			app.ErrorLog.Printf("an error occurred while taking id from url : %v", err)
			http.Error(w, "Taking ID error", http.StatusInternalServerError)
			return
		}

		// take row from database with id
		var habit models.Habit
		query := `SELECT id, title, description, frequency, created_at FROM habits WHERE id = ?`
		row := app.DB.QueryRow(query, id)
		err = row.Scan(&habit.ID, &habit.Title, &habit.Description, &habit.Frequency, &habit.CreatedAt)
		if err != nil {
			app.ErrorLog.Printf("There is no row : %v", err)
			http.Error(w, "There is no row", http.StatusInternalServerError)
			return
		}

		// response
		err = responseSuccess(w, habit, "Habit listed successfully!", http.StatusOK)
		if err != nil {
			app.ErrorLog.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
