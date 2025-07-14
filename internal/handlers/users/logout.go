package users

import (
	"fmt"
	"net/http"
	"time"

	"github.com/furkankorkmaz309/habit-tracker/internal/app"
	"github.com/furkankorkmaz309/habit-tracker/internal/handlers"
)

func Logout(app app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// set cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    "",
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			Expires:  time.Now().Add(-1 * time.Hour),
		})

		// return status
		err := handlers.ResponseSuccess(w, "", "User logged out successfully!", http.StatusOK)
		if err != nil {
			errStr := fmt.Sprintf("an error occurred while encoding json : %v", err)
			http.Error(w, errStr, http.StatusInternalServerError)
			app.ErrorLog.Println(errStr)
			return
		}
	}
}
