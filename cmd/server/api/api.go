package api

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jesusthecreator017/PicoURL/internal/service"
)

const staticDir = "/static"

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

	// Static file serving
	if _, err := os.Stat(staticDir); err == nil {
		// Serve Vite-build assets
		fileServer := http.FileServer(http.Dir(staticDir))
		r.Handle("/assets/*", http.StripPrefix("/", fileServer))

		// Serve the favicon
		r.Get("/favicon.ico", func(w http.ResponseWriter, req *http.Request) {
			http.ServeFile(w, req, filepath.Join(staticDir, "favicon.ico"))
		})

		// Root serves the React SPA
		r.Get("/", func(w http.ResponseWriter, req *http.Request) {
			http.ServeFile(w, req, filepath.Join(staticDir, "index.html"))
		})
	}

	// Redirect route
	r.Get("/{shortcode}", app.handleRedirect)

	return r
}
