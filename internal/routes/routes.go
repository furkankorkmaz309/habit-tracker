package routes

import (
	"net/http"

	"github.com/furkankorkmaz309/habit-tracker/internal/app"
	"github.com/furkankorkmaz309/habit-tracker/internal/handlers"
	"github.com/go-chi/chi"
)

func Routes(app *app.App) http.Handler {
	r := chi.NewRouter()

	r.Route("/habits", func(r chi.Router) {
		r.Post("/", handlers.AddHabit(*app))
		r.Get("/", handlers.GetHabits(*app))
		r.Get("/{id}", handlers.GetHabit(*app))
		r.Delete("/{id}", handlers.DeleteHabit(*app))
		r.Patch("/{id}", handlers.PatchHabit(*app))
	})

	return r
}
