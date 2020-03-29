package handlers

import (
	"net/http"
)

func DebugHandler(ah IAppHandler, w http.ResponseWriter, r *http.Request) (int, error) {

	context := ah.GetContext()
	context.AlignStatus.RLock()
	context.Flags.RLock()
	context.NetworkSettings.APSettings.RLock()
	context.NetworkSettings.RLock()
	context.RLock()

	defer context.RUnlock()
	defer context.NetworkSettings.RUnlock()
	defer context.NetworkSettings.APSettings.RUnlock()
	defer context.AlignStatus.RUnlock()
	defer context.Flags.RUnlock()

	err := writeJson(context, w)
	if err != nil {
		return 500, err
	}

	return 200, nil
}
