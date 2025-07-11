package habits

import (
	"net/http"
	"strconv"

	"github.com/furkankorkmaz309/habit-tracker/internal/app"
	"github.com/go-chi/chi"
)

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
