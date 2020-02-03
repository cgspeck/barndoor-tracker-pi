package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/cgspeck/barndoor-tracker-pi/internal/config"
	"github.com/cgspeck/barndoor-tracker-pi/internal/wireless"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
)

type IAppHandler interface {
	GetContext() *models.AppContext
	GetTime() *time.Time
	// SetTime(*time.Time)
	SetAPMode(bool) error
	GetNetworkSettings() *models.NetworkSettings
}

type AppHandler struct {
	AppContext *models.AppContext
	H          func(IAppHandler, http.ResponseWriter, *http.Request) (int, error)
}

func (ah AppHandler) GetContext() *models.AppContext {
	return ah.AppContext
}

func (ah AppHandler) writeConfig() error {
	// TODO save configuration file
	return nil
}

func (ah *AppHandler) SetAPMode(v bool) error {
	ah.AppContext.NetworkSettings.AccessPointMode = v
	err := wireless.ApplyDesiredConfiguration(ah.GetNetworkSettings())
	if err != nil {
		return err
	}
	err = config.SaveConfig(ah.GetContext())
	return err
}

func (ah AppHandler) GetNetworkSettings() *models.NetworkSettings {
	return ah.AppContext.NetworkSettings
}

func (ah AppHandler) GetTime() *time.Time {
	return ah.AppContext.Time
}

// func (ah AppHandler) SetTime(t *time.Time) {
// 	ah.AppContext.Time = t
// }

func handleHandlerResult(status int, err error, w http.ResponseWriter, r *http.Request) {
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

func (ah AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	status, err := ah.H(&ah, w, r)
	handleHandlerResult(status, err, w, r)
}
