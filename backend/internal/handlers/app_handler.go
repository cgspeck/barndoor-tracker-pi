package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/cgspeck/barndoor-tracker-pi/internal/config"
	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
	"github.com/cgspeck/barndoor-tracker-pi/internal/wireless"
)

type IAppHandler interface {
	GetContext() *models.AppContext
	GetTime() *time.Time
	// SetTime(*time.Time)
	SetAPMode(bool) error
	GetNetworkSettings() *models.NetworkSettings
	GetLocationSettings() *models.LocationSettings
	SetLocationSettings(map[string]interface{}) error

	GetAPSettings() *models.APSettings
	SetAPSettings(map[string]interface{}) error

	GetFlags() *models.Flags
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
	ah.AppContext.Flags.NeedsNetworkSettings = false
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

func (ah AppHandler) GetLocationSettings() *models.LocationSettings {
	return ah.AppContext.LocationSettings
}

func (ah *AppHandler) SetLocationSettings(input map[string]interface{}) error {
	currentSettings := ah.GetLocationSettings()
	mustApplyChanges, newSettings, err := config.IsLocationConfigChanged(input, *currentSettings)
	if err != nil {
		return err
	}

	if ah.GetFlags().NeedsLocationSettings {
		mustApplyChanges = true
	}
	if mustApplyChanges {
		ah.AppContext.LocationSettings = &newSettings
		ah.AppContext.Flags.NeedsLocationSettings = false
		err = config.SaveConfig(ah.GetContext())
		// TODO: hook to recalculate align status
	}
	return err
}

func (ah AppHandler) GetAPSettings() *models.APSettings {
	return ah.AppContext.NetworkSettings.APSettings
}

func (ah *AppHandler) SetAPSettings(input map[string]interface{}) error {
	currentSettings := ah.GetAPSettings()
	mustApplyChanges, newSettings, err := config.IsAPConfigChanged(input, *currentSettings)
	if err != nil {
		return err
	}

	if ah.GetFlags().NeedsNetworkSettings {
		mustApplyChanges = true
	}

	if mustApplyChanges {
		ah.AppContext.NetworkSettings.APSettings = &newSettings
		ah.AppContext.Flags.NeedsNetworkSettings = false
		err = config.SaveConfig(ah.GetContext())
		if err != nil {
			return err
		}
		err = wireless.ApplyDesiredConfiguration(ah.GetNetworkSettings())
	}
	return err
}

func (ah AppHandler) GetTime() *time.Time {
	return ah.AppContext.Time
}

func (ah AppHandler) GetFlags() *models.Flags {
	return ah.AppContext.Flags
}

// func (ah *AppHandler) SetTime(t *time.Time) {
// 	ah.AppContext.Time = t
// }

func writeJson(v interface{}, w http.ResponseWriter) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	io.WriteString(w, string(b))
	return nil
}

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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	status, err := ah.H(&ah, w, r)
	handleHandlerResult(status, err, w, r)
}
