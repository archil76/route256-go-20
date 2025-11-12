package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"route256/cart/internal/infra/logger"
	"strconv"
)

func PrepareID(stringID string) (int64, error) {
	id, err := ConvertID(stringID)
	if err != nil {
		return 0, err
	}

	id, err = ValidateID(id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func ValidateID(id int64) (int64, error) {
	if id < 1 {
		err := errors.New("id should be greater than 0")

		return 0, err
	}

	return id, nil
}

func ConvertID(stringID string) (int64, error) {
	id, err := strconv.ParseInt(stringID, 10, 64)

	if err != nil {

		return 0, err
	}

	return id, nil
}

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

func WriteErrorToResponse(w http.ResponseWriter, r *http.Request, err error, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if len(message) == 0 {
		return
	}

	resp := ErrorResponse{
		Message: message,
	}

	if err != nil {
		resp.Error = err.Error()
	}

	if encodeErr := json.NewEncoder(w).Encode(resp); encodeErr != nil {

		logger.Errorw("failed to write error response",
			r.Method, r.RequestURI, encodeErr, err)
	}
}

func WriteStatusToResponse(w http.ResponseWriter, r *http.Request, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if len(message) == 0 {
		return
	}

	resp := StatusResponse{
		Status: message,
	}

	if encodeErr := json.NewEncoder(w).Encode(resp); encodeErr != nil {
		logger.Errorw("failed to write status response",
			r.Method, r.RequestURI, encodeErr, message)
	}
}

func WriteErrorToLog(r *http.Request, err error, message string) {
	logger.Errorw(message, r.Method, r.RequestURI, err.Error())
}
