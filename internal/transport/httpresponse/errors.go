package httpresponse

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
)

const (
	alreadyExist = "already exist"
)

func Error(w http.ResponseWriter, err error) error {
	if err.Error() == "validation error" {
		return err
	}
	if pgError, ok := err.(*pgconn.PgError); ok {

		if pgError.Code == "23505" {
			return errors.New("already exists")
		}
		return errors.New("server error")

	}
	if _, ok := err.(validator.ValidationErrors); ok {
		return errors.New("validation error")
	}

	if err.Error() == "already exists" {
		return err
	}
	return errors.New("server error")

}
func ReturnError(w http.ResponseWriter, r *http.Request, err error, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(
		Response{
			Status: "error",
			Error:  Error(w, err).Error(),
		},
	)
	return
}
