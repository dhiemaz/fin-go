package common

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type jsonResponse struct {
	Message *string `json:"message,omitempty"` // Message field for the response (optional)
	Status  int     `json:"status"`            // Status code of the response
	Count   *int    `json:"count,omitempty"`   // Count field for the response (optional)
	Data    any     `json:"data,omitempty"`    // Data field for the response (optional)
}

func WriteJSONError(w http.ResponseWriter, message string, statusCode int) {
	writeJSONHeader(w, statusCode)
	response := jsonResponse{
		Message: &message,
		Status:  statusCode,
	}
	encodeJSON(w, response)
}

func WriteJSON(w http.ResponseWriter, statusCode int, data any) {
	writeJSONHeader(w, statusCode)
	response := jsonResponse{
		Status: statusCode,
		Data:   data,
	}
	encodeJSON(w, response)
}

// WriteJSONSimple writes a simple JSON response with the provided status code
func WriteJSONSimple(w http.ResponseWriter, statusCode int, data any) {
	writeJSONHeader(w, statusCode)
	encodeJSON(w, data)
}

// WriteJsonPaginated writes a paginated JSON response to the provided http.ResponseWriter.
// It includes the paginated items, pagination metadata,
// and handles potential errors during pagination or JSON encoding.
func WriteJSONPaginated[T any](w http.ResponseWriter, r *http.Request, items []T, count int64, currentPage int, limit int) {
	writeJSONHeader(w, 200)
	pagination, err := NewPagination(r, items, count, currentPage, limit)
	if err != nil {
		WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := encodeJSON(w, pagination); err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

func writeJSONHeader(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
}

func encodeJSON(w http.ResponseWriter, response interface{}) error {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(response)
	if err != nil {
		return err
	}
	return nil
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}

	return nil
}
