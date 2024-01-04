package healthcheck

import "net/http"

// HandleHealth is a http handler for a simple health  check
func HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
