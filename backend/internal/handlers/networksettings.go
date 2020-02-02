package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
)

func NetworkSettingsHandler(a *models.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	b, err := json.MarshalIndent(a.NetworkSettings, "", "  ")
	if err != nil {
		return 500, err
	}
	io.WriteString(w, string(b))

	return 200, nil
}
