package middlewares

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"net/http"
	"os"

	"github.com/furkankorkmaz309/habit-tracker/internal/app"
	"github.com/furkankorkmaz309/habit-tracker/internal/handlers"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(app app.App) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("token")
			if err == http.ErrNoCookie {
				app.ErrorLog.Printf("there is no cookie : %v", err)
				handlers.ResponseError(w, "there is no cookie", http.StatusInternalServerError)
				return
			}
			if err != nil {
				app.ErrorLog.Printf("an error occurred while reading cookie : %v", err)
				handlers.ResponseError(w, "cookie error", http.StatusInternalServerError)
				return
			}

			// hashed Token
			tokenString := cookie.Value

			// check hashed password
			encryptedData, err := base64.URLEncoding.DecodeString(tokenString)
			if err != nil {
				app.ErrorLog.Printf("base64 decode error: %v", err)
				handlers.ResponseError(w, "Invalid token format", http.StatusUnauthorized)
				return
			}
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

			nonceSize := gcm.NonceSize()
			if len(encryptedData) < nonceSize {
				app.ErrorLog.Printf("an error occurred while comparing hashed string : %v", err)
				handlers.ResponseError(w, "compare error", http.StatusUnauthorized)
				return
			}

			nonce, cipherText := encryptedData[:nonceSize], encryptedData[nonceSize:]
			plaintext, err := gcm.Open(nil, nonce, cipherText, nil)
			if err != nil {
				app.ErrorLog.Printf("decryption error: %v", err)
				handlers.ResponseError(w, "Token decryption failed", http.StatusUnauthorized)
				return
			}

			decryptedJWT := string(plaintext)

			token, err := jwt.ParseWithClaims(decryptedJWT, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("SECRET_KEY")), nil
			})
			if err != nil {
				app.ErrorLog.Printf("an error occurred while parsing jwt : %v", err)
				handlers.ResponseError(w, "JWT parsing error", http.StatusInternalServerError)
				return
			}
			if !token.Valid {
				app.ErrorLog.Printf("token is not valid")
				handlers.ResponseError(w, "Invalid token error", http.StatusUnauthorized)
				return
			}

			claims := token.Claims.(jwt.MapClaims)
			userIdFloat, ok := claims["userId"].(float64)
			if !ok {
				app.ErrorLog.Printf("Unauthorized")
				handlers.ResponseError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			userId := int(userIdFloat)

			username, ok := claims["username"].(string)
			if !ok {
				app.ErrorLog.Printf("Unauthorized")
				handlers.ResponseError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "userId", userId)
			ctx = context.WithValue(ctx, "username", username)

			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
