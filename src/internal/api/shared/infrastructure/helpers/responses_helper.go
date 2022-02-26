package helpers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// WriteOkResponse is used when and endpoind does not respond with an error
func WriteOkResponse(r *http.Request, w http.ResponseWriter, statusCode int, content interface{}) {
	log.Printf("[%v] %v %v", GetRequestIDFromContext(r), statusCode, time.Since(GetRequestStartTimeFromContext(r)))

	if content == nil {
		w.WriteHeader(statusCode)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(content)
}

// WriteErrorResponse is used when and endpoind responds with an error
func WriteErrorResponse(r *http.Request, w http.ResponseWriter, statusCode int, msg string, internalError error) {
	requestID := GetRequestIDFromContext(r)

	if internalError != nil {
		log.Printf("[%v] %v %v %v (%v)", requestID, statusCode, time.Since(GetRequestStartTimeFromContext(r)), msg, internalError)
	} else {
		log.Printf("[%v] %v %v %v", requestID, statusCode, time.Since(GetRequestStartTimeFromContext(r)), msg)
	}

	http.Error(w, msg, statusCode)
}
