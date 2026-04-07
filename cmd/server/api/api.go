package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jesusthecreator017/PicoURL/internal/service"
)

type Application struct {
	service service.URLService
	router  *chi.Mux
}

func NewApplication(svc service.URLService) *Application {
	app := &Application{service: svc}
	app.router = app.setupRoutes()
	return app
}

// Expose the router as an http.Handler so main file can use it
func (app *Application) Handler() http.Handler {
	return app.router
}

func (app *Application) setupRoutes() *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// API routes
	r.Post("/api/shorten", app.handleShorten)
	r.Get("/api/stats/{shortcode}", app.handleStats)
	r.Delete("/api/{shortcode}", app.handleDelete)

	// Redirect route
	r.Get("/{shortcode}", app.handleRedirect)

	return r
}
