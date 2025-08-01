package users

import (
	"fmt"
	"net/http"

	"github.com/furkankorkmaz309/habit-tracker/internal/app"
	"github.com/furkankorkmaz309/habit-tracker/internal/handlers"
)

func Logout(app app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// set cookie
		http.SetCookie(w, createCookie("token", "", -1))
		http.SetCookie(w, createCookie("uid", "", -1))

		// return status
		err := handlers.ResponseSuccess(w, "", "User logged out successfully!", http.StatusOK)
		if err != nil {
			errStr := fmt.Sprintf("an error occurred while encoding json : %v", err)
			handlers.ResponseError(w, errStr, http.StatusInternalServerError)
			app.ErrorLog.Println(errStr)
			return
		}
	}
}
