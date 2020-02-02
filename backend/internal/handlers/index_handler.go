package handlers

import (
	"fmt"
	"io"
	"net/http"
)

func IndexHandler(ah IAppHandler, w http.ResponseWriter, r *http.Request) (int, error) {
	io.WriteString(w, fmt.Sprintf("%v", ah.GetTime()))
	return 200, nil
}
