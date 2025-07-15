package habits

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/furkankorkmaz309/habit-tracker/internal/app"
	"github.com/furkankorkmaz309/habit-tracker/internal/handlers"
	"github.com/furkankorkmaz309/habit-tracker/internal/models"
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
		fq, err := handlers.FrequencyConvert(habit.Frequency)
		if err != nil {
			app.ErrorLog.Printf("an error occurred while converting frequency : %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// add habit to database
		habit.CreatedAt = time.Now()
		habit.UserID = r.Context().Value("userId").(int)
		query := `INSERT INTO habits (title, description, frequency, created_at, user_id) VALUES (?, ?, ?, ?, ?)`
		result, err := app.DB.Exec(query, habit.Title, habit.Description, fq, habit.CreatedAt, habit.UserID)
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
		err = handlers.ResponseSuccess(w, habit, "Habit added successfully!", http.StatusCreated)
		if err != nil {
			app.ErrorLog.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
