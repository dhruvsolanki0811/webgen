package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/dhruvsolanki0811/webgen/internal/domain"
)

func RespondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func RespondError(w http.ResponseWriter, err error) {
	var appErr *domain.AppError
	if !errors.As(err, &appErr) {
		appErr = domain.ErrInternal(err)
	}

	if appErr.Err != nil {
		log.Printf("ERROR [%d]: %s — %v", appErr.Code, appErr.Message, appErr.Err)
	}

	RespondJSON(w, appErr.Code, map[string]string{"error": appErr.Message})
}

func Decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, domain.ErrBadRequest("invalid request body")
	}
	return v, nil
}
