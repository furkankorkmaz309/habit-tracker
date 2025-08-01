package habits

import (
	"encoding/json"
	"net/http"
	"strings"
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
			handlers.ResponseError(w, "Decode error", http.StatusInternalServerError)
			return
		}
		if strings.TrimSpace(habit.Title) == "" {
			app.ErrorLog.Printf("Title can not blank")
			handlers.ResponseError(w, "Title can not blank", http.StatusInternalServerError)
			return
		}
		if strings.TrimSpace(habit.Description) == "" {
			app.ErrorLog.Printf("Description can not blank")
			handlers.ResponseError(w, "Description can not blank", http.StatusInternalServerError)
			return
		}
		if strings.TrimSpace(habit.Frequency) == "" {
			app.ErrorLog.Printf("Frequency can not blank")
			handlers.ResponseError(w, "Frequency can not blank", http.StatusInternalServerError)
			return
		}

		if habit.Frequency != "D" && habit.Frequency != "3D" && habit.Frequency != "W" && habit.Frequency != "2W" && habit.Frequency != "3W" && habit.Frequency != "M" {
			app.ErrorLog.Println("Only D, 3D, W, 2W, 3W, M")
			handlers.ResponseError(w, "Only D, 3D, W, 2W, 3W, M", http.StatusInternalServerError)
			return
		}

		if strings.TrimSpace(habit.Day) == "" && habit.Frequency != "D" {
			app.ErrorLog.Printf("Day can not blank")
			handlers.ResponseError(w, "Day can not blank", http.StatusInternalServerError)
			return
		}

		if habit.Frequency != "D" && habit.Day != "monday" && habit.Day != "tuesday" && habit.Day != "wednesday" && habit.Day != "thursday" && habit.Day != "friday" && habit.Day != "saturday" && habit.Day != "sunday" {
			errStr := `Only 'monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday', 'sunday'`
			app.ErrorLog.Println(errStr)
			handlers.ResponseError(w, errStr, http.StatusInternalServerError)
			return
		}

		// add habit to database
		habit.CreatedAt = time.Now()
		habit.UserID = r.Context().Value("userId").(int)
		query := `INSERT INTO habits (title, description, frequency, day, created_at, user_id) VALUES (?, ?, ?, ?, ?, ?)`
		result, err := app.DB.Exec(query, habit.Title, habit.Description, habit.Frequency, habit.Day, habit.CreatedAt, habit.UserID)
		if err != nil {
			app.ErrorLog.Printf("an error occurred while inserting habit to database : %v", err)
			handlers.ResponseError(w, "DB insert error", http.StatusInternalServerError)
			return
		}

		id64, err := result.LastInsertId()
		if err != nil {
			app.ErrorLog.Printf("an error occurred while retrieve ID : %v", err)
			handlers.ResponseError(w, "ID retrieve error", http.StatusInternalServerError)
			return
		}
		habit.ID = int(id64)

		// response habit and status code
		err = handlers.ResponseSuccess(w, habit, "Habit added successfully!", http.StatusCreated)
		if err != nil {
			app.ErrorLog.Println(err)
			handlers.ResponseError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
