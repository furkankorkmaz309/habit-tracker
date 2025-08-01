package users

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
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

		// hash token
		encKey := os.Getenv("ENCRYPTION_KEY")
		if len(encKey) != 32 {
			app.ErrorLog.Printf("ENCRYPTION_KEY must be 32 bytes, got %d", len(encKey))
			handlers.ResponseError(w, "Invalid encryption key", http.StatusInternalServerError)
			return
		}
		if encKey == "" {
			app.ErrorLog.Printf("there is no ENCRYPTION_KEY in .env file")
			handlers.ResponseError(w, ".env error", http.StatusInternalServerError)
			return
		}
		c, err := aes.NewCipher([]byte(encKey))
		if err != nil {
			app.ErrorLog.Printf("an error occurred while generating cipher : %v", err)
			handlers.ResponseError(w, "Cipher generate error", http.StatusUnauthorized)
			return
		}

		gcm, err := cipher.NewGCM(c)
		if err != nil {
			app.ErrorLog.Printf("an error occurred while generating GCM : %v", err)
			handlers.ResponseError(w, "GCM generate error", http.StatusUnauthorized)
			return
		}
		nonce := make([]byte, gcm.NonceSize())
		rand.Read(nonce)

		encrypted := gcm.Seal(nonce, nonce, []byte(signedString), nil)
		encryptedJWT := base64.URLEncoding.EncodeToString(encrypted)

		nonceID := make([]byte, gcm.NonceSize())
		rand.Read(nonceID)
		userIDBytes := []byte(fmt.Sprintf("%d", user.ID))
		encryptedUserID := gcm.Seal(nonceID, nonceID, userIDBytes, nil)
		encryptedUserIDString := base64.URLEncoding.EncodeToString(encryptedUserID)

		// set cookie
		http.SetCookie(w, createCookie("token", encryptedJWT, 3600))
		http.SetCookie(w, createCookie("uid", encryptedUserIDString, 3600))

		user.ID = 0
		user.Password = ""
		// return status
		err = handlers.ResponseSuccess(w, user, "User authenticated successfully!", http.StatusOK)
		// err = handlers.ResponseLogin(w, user, "User authenticated successfully!", signedString, http.StatusOK)
		if err != nil {
			errStr := fmt.Sprintf("an error occurred while encoding json : %v", err)
			handlers.ResponseError(w, errStr, http.StatusInternalServerError)
			app.ErrorLog.Println(errStr)
			return
		}
	}
}
