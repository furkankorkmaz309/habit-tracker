package middlewares

import (
	"context"
	"net/http"
	"os"

	"github.com/furkankorkmaz309/habit-tracker/internal/app"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(app app.App) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("token")
			if err == http.ErrNoCookie {
				app.ErrorLog.Printf("there is no cookie : %v", err)
				http.Error(w, "there is no cookie", http.StatusInternalServerError)
				return
			}
			if err != nil {
				app.ErrorLog.Printf("an error occurred while reading cookie : %v", err)
				http.Error(w, "cookie error", http.StatusInternalServerError)
				return
			}

			tokenString := cookie.Value
			token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("SECRET_KEY")), nil
			})
			if err != nil {
				app.ErrorLog.Printf("an error occurred while parsing jwt : %v", err)
				http.Error(w, "JWT parsing error", http.StatusInternalServerError)
				return
			}
			if !token.Valid {
				app.ErrorLog.Printf("token is not valid")
				http.Error(w, "Invalid token error", http.StatusUnauthorized)
				return
			}

			claims := token.Claims.(jwt.MapClaims)
			userIdFloat, ok := claims["userId"].(float64)
			if !ok {
				app.ErrorLog.Printf("Unauthorized")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			userId := int(userIdFloat)

			username, ok := claims["username"].(string)
			if !ok {
				app.ErrorLog.Printf("Unauthorized")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "userId", userId)
			ctx = context.WithValue(ctx, "username", username)

			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
