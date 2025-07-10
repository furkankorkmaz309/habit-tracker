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

func DeleteHabit(app app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// take id from url
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			app.ErrorLog.Printf("an error occurred while taking id from url : %v", err)
			http.Error(w, "Taking ID error", http.StatusInternalServerError)
			return
		}

		// delete from database
		query := `DELETE FROM habits WHERE id = ?`
		result, err := app.DB.Exec(query, id)
		if err != nil {
			app.ErrorLog.Printf("an error occurred while deleting habit : %v", err)
			http.Error(w, "Delete error", http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			app.ErrorLog.Printf("an error occurred while retrieving result : %v", err)
			http.Error(w, "Retrieve error", http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			app.ErrorLog.Printf("there is no row with id : %v", id)
			http.Error(w, "No row with that id", http.StatusNotFound)
			return
		}

		// response
		w.WriteHeader(http.StatusNoContent)
	}
}

func PatchHabit(app app.App) http.HandlerFunc {
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

		// pick the parts will change
		var habitInput models.Habit
		err = json.NewDecoder(r.Body).Decode(&habitInput)
		if err != nil {
			app.ErrorLog.Printf("an error occurred while decoding input : %v", err)
			http.Error(w, "Decode error", http.StatusInternalServerError)
			return
		}

		if habitInput.Title != "" {
			habit.Title = habitInput.Title
		}

		if habitInput.Description != "" {
			habit.Description = habitInput.Description
		}

		fq, err := strconv.Atoi(habit.Frequency)
		if err != nil {
			app.ErrorLog.Printf("an error occurred while converting string to integer : %v", err)
			http.Error(w, "Conversion error", http.StatusInternalServerError)
			return
		}
		if habitInput.Frequency != "" {
			fq, err = FrequencyConvert(habitInput.Frequency)
			if err != nil {
				app.ErrorLog.Printf("an error occurred while converting frequency : %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		// change database
		queryUpdate := `UPDATE habits SET title = ?, description = ?, frequency = ? WHERE id = ?`
		_, err = app.DB.Exec(queryUpdate, habit.Title, habit.Description, fq, id)
		if err != nil {
			app.ErrorLog.Printf("an error occurred while updating habit : %v", err)
			http.Error(w, "Update error", http.StatusInternalServerError)
			return
		}

		// response
		err = responseSuccess(w, habit, "Habit updated successfully!", http.StatusOK)
		if err != nil {
			app.ErrorLog.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}
