package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status   string      `json:"status"`
	Error    string      `json:"error,omitempty"`
	Response interface{} `json:"response,omitempty"`
}

func ReturnOk(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(
		Response{
			Status: "ok",
		},
	)
	return
}
func ReturnOkData(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(
		Response{
			Status:   "ok",
			Response: data,
		},
	)
	return
}
