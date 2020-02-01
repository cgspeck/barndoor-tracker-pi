package handlers

import (
	"log"
	"net/http"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
)

type AppHandler struct {
	*models.AppContext
	H func(*models.AppContext, http.ResponseWriter, *http.Request) (int, error)
}

func (ah AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, err := ah.H(ah.AppContext, w, r)

	if err != nil {
		log.Printf("HTTP %d: %q", status, err)
		switch status {
		case http.StatusNotFound:
			http.NotFound(w, r)
			// And if we wanted a friendlier error page, we can
			// now leverage our context instance - e.g.
			// err := ah.renderTemplate(w, "http_404.tmpl", nil)
		case http.StatusInternalServerError:
			http.Error(w, http.StatusText(status), status)
		default:
			http.Error(w, http.StatusText(status), status)
		}
	}
	w.Header().Set("Content-Type", "application/json")
}
