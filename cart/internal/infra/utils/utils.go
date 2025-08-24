package utils

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func PrepareID(w http.ResponseWriter, r *http.Request, stringID string) (int64, error) {

	id, err := ConvertID(w, r, stringID)
	if err != nil {
		return id, err
	}

	id, err = ValidateID(w, r, id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func ValidateID(w http.ResponseWriter, r *http.Request, id int64) (int64, error) {

	if id < 1 {

		err := errors.New("id should be greater than 0")
		err = WriteErrorToResponse(w, r, err, "", http.StatusBadRequest)

		return id, err
	}
	return id, nil
}

func ConvertID(w http.ResponseWriter, r *http.Request, stringID string) (int64, error) {

	id, err := strconv.ParseInt(stringID, 10, 64)

	if err != nil {
		err = WriteErrorToResponse(w, r, err, "parsing error", http.StatusBadRequest)

		return id, err
	}
	return id, nil
}

func WriteErrorToResponse(w http.ResponseWriter, r *http.Request, err error, message string, status int) error {

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	_, errOut := fmt.Fprintf(w, "{\"message\": \"%s\" - \"%s\"}", message, err)
	if errOut != nil {
		//instead r.pat.str
		log.Printf("%s %s: %s - %s", r.Method, r.RequestURI, errOut.Error(), message)
		return errOut
	}

	return err

}

func WriteStatusToResponse(w http.ResponseWriter, r *http.Request, message string, status int) {

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	_, errOut := fmt.Fprintf(w, "{%s}", message)
	if errOut != nil {
		log.Printf("%s %s: %s - %s", r.Method, r.RequestURI, errOut.Error(), message)
	}

}
