package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/furkankorkmaz309/habit-tracker/internal/app"
	"github.com/furkankorkmaz309/habit-tracker/internal/handlers"
	"github.com/furkankorkmaz309/habit-tracker/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Login(app app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// take user input
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			errStr := fmt.Sprintf("an error occurred while decoding json : %v", err)
			handlers.ResponseError(w, errStr, http.StatusBadRequest)
			app.ErrorLog.Println(errStr)
			return
		}

		// check is data exists username and if user exists take passwordv
		var count int
		var id int
		var email string
		var password string
		var createdAt time.Time
		queryUserCheck := `SELECT COUNT(*), id, email, password, created_at FROM users WHERE username = ?`
		err = app.DB.QueryRow(queryUserCheck, user.Username).Scan(&count, &id, &email, &password, &createdAt)
		if count == 0 {
			errStr := fmt.Sprintf("there is no user with username : %v", user.Username)
			handlers.ResponseError(w, errStr, http.StatusInternalServerError)
			app.ErrorLog.Println(errStr)
			return
		}
		if err != nil {
			errStr := fmt.Sprintf("an error occurred while checking is user available : %v", err)
			handlers.ResponseError(w, errStr, http.StatusInternalServerError)
			app.ErrorLog.Println(errStr)
			return
		}

		// check hashed password
		err = bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password))
		if err != nil {
			errStr := "wrong username or password"
			handlers.ResponseError(w, errStr, http.StatusInternalServerError)
			app.ErrorLog.Println(errStr)
			return
		}
		user.ID = id
		user.CreatedAt = createdAt
		user.Email = email

		// jwt
		key := os.Getenv("SECRET_KEY")
		if key == "" {
			app.ErrorLog.Printf("there is no SECRET_KEY in .env file")
			handlers.ResponseError(w, ".env error", http.StatusInternalServerError)
			return
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userId":   user.ID,
			"username": user.Username,
			"iat":      time.Now().Unix(),
			"exp":      time.Now().Add(1 * time.Hour).Unix(),
		})

		signedString, err := token.SignedString([]byte(key))
		if err != nil {
			app.ErrorLog.Printf("an error occurred while signing string : %v", err)
			handlers.ResponseError(w, "Could not authanticate user", http.StatusUnauthorized)
			return
		}

		// set cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    signedString,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			Expires:  time.Now().Add(1 * time.Hour),
		})

		// return status
		err = handlers.ResponseSuccess(w, user, "User authenticated successfully!", http.StatusOK)
		if err != nil {
			errStr := fmt.Sprintf("an error occurred while encoding json : %v", err)
			handlers.ResponseError(w, errStr, http.StatusInternalServerError)
			app.ErrorLog.Println(errStr)
			return
		}
	}
}
