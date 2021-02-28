package helpers

import (
	"encoding/json"
	"log"
	"net/http"
)

// WriteOkResponse is used when and endpoind does not respond with an error
func WriteOkResponse(r *http.Request, w http.ResponseWriter, statusCode int, content interface{}) {
	log.Printf("[%v] %v", GetRequestIDFromContext(r), statusCode)

	const jsonContentType = "application/json"

	if content != nil {
		w.Header().Set("content-type", jsonContentType)
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(content)
	} else {
		w.WriteHeader(statusCode)
	}
}

// WriteErrorResponse is used when and endpoind responds with an error
func WriteErrorResponse(r *http.Request, w http.ResponseWriter, statusCode int, msg string, internalError error) {
	requestID := GetRequestIDFromContext(r)
	if internalError != nil {
		log.Printf("[%v] %v %v (%v)", requestID, statusCode, msg, internalError)
	} else {
		log.Printf("[%v] %v %v", requestID, statusCode, msg)
	}
	http.Error(w, msg, statusCode)
}
