package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

func DebugHandler(ah IAppHandler, w http.ResponseWriter, r *http.Request) (int, error) {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent(ah.GetContext(), "", "  ")
	if err != nil {
		return 500, err
	}
	io.WriteString(w, string(b))

	return 200, nil
}
