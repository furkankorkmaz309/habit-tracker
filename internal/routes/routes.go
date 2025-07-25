package routes

import (
	"net/http"

	"github.com/furkankorkmaz309/habit-tracker/internal/app"
	"github.com/furkankorkmaz309/habit-tracker/internal/handlers/habits"
	"github.com/furkankorkmaz309/habit-tracker/internal/handlers/middlewares"
	"github.com/furkankorkmaz309/habit-tracker/internal/handlers/users"
	"github.com/go-chi/chi"
)

func Routes(app *app.App) http.Handler {
	r := chi.NewRouter()

	r.Use(middlewares.CORS)

	r.Route("/habits", func(r chi.Router) {
		r.Use(middlewares.Auth(*app))

		r.Post("/", habits.AddHabit(*app))
		r.Get("/", habits.GetHabits(*app))
		r.Get("/{id}", habits.GetHabit(*app))
		r.Delete("/{id}", habits.DeleteHabit(*app))
		r.Patch("/{id}", habits.PatchHabit(*app))
	})

	r.Post("/signup", users.Signup(*app))
	r.Post("/login", users.Login(*app))
	r.Post("/logout", users.Logout(*app))

	return r
}
