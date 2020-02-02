package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
)

type IAppHandler interface {
	GetContext() *models.AppContext
	GetTime() *time.Time
	// SetTime(*time.Time)
	WriteConfig()
	SetAPMode(bool)
	GetNetworkSettings() *models.NetworkSettings
}

type AppHandler struct {
	AppContext *models.AppContext
	H          func(IAppHandler, http.ResponseWriter, *http.Request) (int, error)
}

func (ah AppHandler) GetContext() *models.AppContext {
	return ah.AppContext
}

func (ah AppHandler) WriteConfig()   {}
func (ah AppHandler) SetAPMode(bool) {}

func (ah AppHandler) GetNetworkSettings() *models.NetworkSettings {
	return ah.AppContext.NetworkSettings
}

func (ah AppHandler) GetTime() *time.Time {
	return ah.AppContext.Time
}

// func (ah AppHandler) SetTime(t *time.Time) {
// 	ah.AppContext.Time = t
// }

func (ah AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, err := ah.H(&ah, w, r)

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
}
