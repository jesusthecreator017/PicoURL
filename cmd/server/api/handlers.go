package api

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jesusthecreator017/PicoURL/cmd/server/api/helpers"
	"github.com/jesusthecreator017/PicoURL/internal/service"
)

func (app *Application) handleShorten(w http.ResponseWriter, req *http.Request) {
	// Load the request in a variable
	var r shortenRequest
	if err := helpers.ReadJson(req, &r); err != nil {
		helpers.WriteJson(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}

	// Check each field in the request
	if r.URL == "" {
		helpers.WriteJson(w, http.StatusBadRequest, errorResponse{Error: "url is required"})
		return
	}

	// Perform the shortening
	shortCode, err := app.service.Shorten(req.Context(), r.URL)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidURL):
			helpers.WriteJson(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		case errors.Is(err, service.ErrUnreachableURL):
			helpers.WriteJson(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		case errors.Is(err, service.ErrCollision):
			helpers.WriteJson(w, http.StatusInternalServerError, errorResponse{Error: err.Error()})
		default:
			log.Printf("shorten error: %v", err)
			helpers.WriteJson(w, http.StatusInternalServerError, errorResponse{Error: "internal server error"})
		}
		return
	}

	helpers.WriteJson(w, http.StatusCreated, shortenResponse{ShortURL: shortCode})
}

func (app *Application) handleRedirect(w http.ResponseWriter, req *http.Request) {
	// Exctract the shortcode from url params
	shortCode := chi.URLParam(req, "shortcode")

	originalURL, err := app.service.Resolve(req.Context(), shortCode)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			helpers.WriteJson(w, http.StatusNotFound, errorResponse{Error: "short URL not found"})
			return
		}
		helpers.WriteJson(w, http.StatusInternalServerError, errorResponse{Error: "internal server error"})
		return
	}

	http.Redirect(w, req, originalURL, http.StatusMovedPermanently)
}

func (app *Application) handleStats(w http.ResponseWriter, req *http.Request) {
	// Read path params
	shortCode := chi.URLParam(req, "shortcode")

	// Get original stats
	count, err := app.service.GetStats(req.Context(), shortCode)
	if err != nil {
		helpers.WriteJson(w, http.StatusNotFound, errorResponse{Error: "short URL not found"})
		return
	}

	helpers.WriteJson(w, http.StatusOK, statsResponse{
		ShortURL:   shortCode,
		ClickCount: count,
	})
}

func (app *Application) handleDelete(w http.ResponseWriter, req *http.Request) {
	// Read pathparams
	shortCode := chi.URLParam(req, "shortcode")

	if err := app.service.Delete(req.Context(), shortCode); err != nil {
		helpers.WriteJson(w, http.StatusInternalServerError, errorResponse{Error: "failed to delete"})
		return
	}

	helpers.WriteJson(w, http.StatusOK, helpers.Envelope{"message": "deleted"})
}
