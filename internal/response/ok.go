package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status string      `json:"status"`
	Error  string      `json:"error,omitempty"`
	Body   interface{} `json:"body,omitempty"`
}

func ReturnOk(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(
		Response{
			Status: "ok",
		},
	)
	return
}
func ReturnOkData(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(
		Response{
			Status: "ok",
			Body:   data,
		},
	)
	return
}
