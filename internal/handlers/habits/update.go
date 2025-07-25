package habits

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/furkankorkmaz309/habit-tracker/internal/app"
	"github.com/furkankorkmaz309/habit-tracker/internal/handlers"
	"github.com/furkankorkmaz309/habit-tracker/internal/models"
	"github.com/go-chi/chi"
)

func PatchHabit(app app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// take id from url
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			app.ErrorLog.Printf("an error occurred while taking id from url : %v", err)
			handlers.ResponseError(w, "Taking ID error", http.StatusInternalServerError)
			return
		}

		// take row from database with id
		var habit models.Habit
		habit.UserID = r.Context().Value("userId").(int)
		query := `SELECT id, title, description, frequency, created_at FROM habits WHERE id = ? AND user_id = ?`
		row := app.DB.QueryRow(query, id, habit.UserID)
		err = row.Scan(&habit.ID, &habit.Title, &habit.Description, &habit.Frequency, &habit.CreatedAt)
		if err != nil {
			app.ErrorLog.Printf("There is no row : %v", err)
			handlers.ResponseError(w, "There is no row", http.StatusInternalServerError)
			return
		}

		// pick the parts will change
		var habitInput models.Habit
		err = json.NewDecoder(r.Body).Decode(&habitInput)
		if err != nil {
			app.ErrorLog.Printf("an error occurred while decoding input : %v", err)
			handlers.ResponseError(w, "Decode error", http.StatusInternalServerError)
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
			handlers.ResponseError(w, "Conversion error", http.StatusInternalServerError)
			return
		}
		if habitInput.Frequency != "" {
			fq, err = handlers.FrequencyConvert(habitInput.Frequency)
			if err != nil {
				app.ErrorLog.Printf("an error occurred while converting frequency : %v", err)
				handlers.ResponseError(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		// change database
		queryUpdate := `UPDATE habits SET title = ?, description = ?, frequency = ? WHERE id = ?`
		_, err = app.DB.Exec(queryUpdate, habit.Title, habit.Description, fq, id)
		if err != nil {
			app.ErrorLog.Printf("an error occurred while updating habit : %v", err)
			handlers.ResponseError(w, "Update error", http.StatusInternalServerError)
			return
		}

		// response
		err = handlers.ResponseSuccess(w, habit, "Habit updated successfully!", http.StatusOK)
		if err != nil {
			app.ErrorLog.Println(err)
			handlers.ResponseError(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}
