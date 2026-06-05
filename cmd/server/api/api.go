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
	service    service.URLService
	router     *chi.Mux
	corsOrigin string
}

func NewApplication(svc service.URLService, corsOrigin string) *Application {
	app := &Application{service: svc, corsOrigin: corsOrigin}
	app.router = app.setupRoutes()
	return app
}

// Expose the router as an http.Handler so main file can use it
func (app *Application) Handler() http.Handler {
	return app.router
}

// corsMiddleware allows the API to be called from a browser on another origin
// (e.g. the personal site embedding the live demo). origin is the allowed origin
// ("*" by default); browser preflight requests are answered with 204.
func corsMiddleware(origin string) func(http.Handler) http.Handler {
	if origin == "" {
		origin = "*"
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Vary", "Origin")
			if req.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, req)
		})
	}
}

func (app *Application) setupRoutes() *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware(app.corsOrigin))

	// API routes
	r.Post("/api/shorten", app.handleShorten)
	r.Get("/api/stats/{shortcode}", app.handleStats)
	r.Delete("/api/{shortcode}", app.handleDelete)
	r.Get("/api/total", app.handleTotal)

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
