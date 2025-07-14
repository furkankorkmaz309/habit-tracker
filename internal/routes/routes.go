package routes

import (
	"net/http"

	"github.com/furkankorkmaz309/habit-tracker/internal/app"
	"github.com/furkankorkmaz309/habit-tracker/internal/handlers/habits"
	"github.com/furkankorkmaz309/habit-tracker/internal/handlers/users"
	"github.com/go-chi/chi"
)

func Routes(app *app.App) http.Handler {
	r := chi.NewRouter()

	r.Route("/habits", func(r chi.Router) {
		r.Post("/", habits.AddHabit(*app))
		r.Get("/", habits.GetHabits(*app))
		r.Get("/{id}", habits.GetHabit(*app))
		r.Delete("/{id}", habits.DeleteHabit(*app))
		r.Patch("/{id}", habits.PatchHabit(*app))
	})

	r.Post("/signup", users.Signup(*app))
	r.Post("/login", users.Login(*app))

	return r
}
