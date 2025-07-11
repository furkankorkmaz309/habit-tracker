package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/furkankorkmaz309/habit-tracker/internal/app"
	"github.com/furkankorkmaz309/habit-tracker/internal/handlers"
	"github.com/furkankorkmaz309/habit-tracker/internal/models"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

func Signup(app app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// take user input
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			errStr := fmt.Sprintf("an error occurred while decoding json : %v", err)
			http.Error(w, errStr, http.StatusBadRequest)
			app.ErrorLog.Println(errStr)
			return
		}

		var usernameValidationErrors []string
		var usernameMaxLen string
		var usernameMinLen string

		var passwordValidationErrors []string
		var passwordMaxLen string
		var passwordMinLen string
		var passwordNumber string
		var passwordUppercase string
		var passwordLowercase string

		// check are username and password correct type?
		if len(user.Username) > 20 {
			usernameMaxLen = "- Must be at most 20 characters"
			usernameValidationErrors = append(usernameValidationErrors, usernameMaxLen)
		}
		if len(user.Username) < 8 {
			usernameMinLen = "- Must be at least 8 characters"
			usernameValidationErrors = append(usernameValidationErrors, usernameMinLen)
		}

		err = validator.New().Struct(user)
		if err != nil {
			passwordNumber = "- Must contain at least one number"
			passwordValidationErrors = append(passwordValidationErrors, passwordNumber)
		}
		if user.Password == strings.ToUpper(user.Password) {
			passwordUppercase = "- Must contain at least one uppercase letter"
			passwordValidationErrors = append(passwordValidationErrors, passwordUppercase)
		}
		if user.Password == strings.ToLower(user.Password) {
			passwordLowercase = "- Must contain at least one lowercase letter"
			passwordValidationErrors = append(passwordValidationErrors, passwordLowercase)
		}
		if len(user.Password) > 25 {
			passwordMaxLen = "- Must be at most 25 characters"
			passwordValidationErrors = append(passwordValidationErrors, passwordMaxLen)
		}
		if len(user.Password) < 12 {
			passwordMinLen = "- Must be at least 12 characters"
			passwordValidationErrors = append(passwordValidationErrors, passwordMinLen)
		}

		if len(usernameValidationErrors) > 0 && len(passwordValidationErrors) > 0 {
			errStrUser := fmt.Sprint("Username is invalid\n" + strings.Join(usernameValidationErrors, "\n"))
			errStrPassword := fmt.Sprint("Password is invalid\n" + strings.Join(passwordValidationErrors, "\n"))

			app.ErrorLog.Println(errStrUser)
			app.ErrorLog.Println(errStrPassword)

			http.Error(w, errStrUser+"\n"+errStrPassword, http.StatusBadRequest)
			return
		}

		if len(usernameValidationErrors) > 0 {
			errStrUser := fmt.Sprint("Username is invalid\n" + strings.Join(usernameValidationErrors, "\n"))
			app.ErrorLog.Println(errStrUser)
			http.Error(w, errStrUser, http.StatusBadRequest)
			return
		}
		if len(passwordValidationErrors) > 0 {
			errStrPassword := fmt.Sprint("Password is invalid\n" + strings.Join(passwordValidationErrors, "\n"))
			app.ErrorLog.Println(errStrPassword)
			http.Error(w, errStrPassword, http.StatusBadRequest)
			return
		}

		// check is username available
		var count int
		queryUserCheck := `SELECT COUNT(*) FROM users WHERE username = ?`
		err = app.DB.QueryRow(queryUserCheck, user.Username).Scan(&count)
		if err != nil {
			errStr := fmt.Sprintf("an error occurred while checking is user available : %v", err)
			app.ErrorLog.Println(errStr)
			http.Error(w, errStr, http.StatusInternalServerError)
			return
		}

		if count > 0 {
			errStr := "Username already taken"
			app.ErrorLog.Println(errStr)
			http.Error(w, errStr, http.StatusBadRequest)
			return
		}

		// take cost from .env
		costStr := os.Getenv("COST_NUM")
		if costStr == "" {
			app.ErrorLog.Println("COST_NUM env variable not set")
			http.Error(w, "Server configuration error", http.StatusInternalServerError)
			return
		}
		cost, err := strconv.Atoi(costStr)
		if err != nil {
			app.ErrorLog.Printf("an error occurred while converting string to integer : %v", err)
			http.Error(w, "Conversion error", http.StatusInternalServerError)
			return
		}

		// hash password
		hashedPasswordByte, err := bcrypt.GenerateFromPassword([]byte(user.Password), cost)
		if err != nil {
			errStr := fmt.Sprintf("an error occurred while hashing password : %v", err)
			http.Error(w, errStr, http.StatusInternalServerError)
			app.ErrorLog.Println(errStr)
			return
		}
		hashedPassword := string(hashedPasswordByte)

		user.Password = hashedPassword
		user.CreatedAt = time.Now()

		// add to database
		query := `INSERT INTO users (username, password, created_at) VALUES (?,?,?)`
		result, err := app.DB.Exec(query, user.Username, user.Password, user.CreatedAt)
		if err != nil {
			errStr := fmt.Sprintf("an error occurred while inserting user to database : %v", err)
			http.Error(w, errStr, http.StatusInternalServerError)
			app.ErrorLog.Println(errStr)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			errStr := fmt.Sprintf("an error occurred while retrieving inserted ID : %v", err)
			http.Error(w, errStr, http.StatusInternalServerError)
			app.ErrorLog.Println(errStr)
			return
		}
		user.ID = int(id)

		// return status
		err = handlers.ResponseSuccess(w, user, "User authenticated successfully!", http.StatusCreated)
		if err != nil {
			errStr := fmt.Sprintf("an error occurred while encoding json : %v", err)
			http.Error(w, errStr, http.StatusInternalServerError)
			app.ErrorLog.Println(errStr)
			return
		}

	}
}
