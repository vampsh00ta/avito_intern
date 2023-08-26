package response

import (
	"encoding/json"
	"net/http"
)

const (
	alreadyExist = "already exist"
)

func whichError(w http.ResponseWriter, err error) {
	switch err.Error() {
	case alreadyExist:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
func ReturnError(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(
		Response{
			Status: "error",
			Error:  err.Error(),
		},
	)
	return
}
