package helpers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/consts"
)

// WriteOkResponse is used when and endpoind does not respond with an error
func WriteOkResponse(r *http.Request, w http.ResponseWriter, statusCode int, content interface{}) {
	log.Printf("[%v] %v %v", GetRequestIDFromContext(r), statusCode, time.Since(getRequestStartTimeFromContext(r)))

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
	timeSinceReqStart := time.Since(getRequestStartTimeFromContext(r))

	if internalError != nil {
		log.Printf("[%v] %v %v %v (%v)", requestID, statusCode, timeSinceReqStart, msg, internalError)
	} else {
		log.Printf("[%v] %v %v %v", requestID, statusCode, timeSinceReqStart, msg)
	}

	http.Error(w, msg, statusCode)
}

func getRequestStartTimeFromContext(r *http.Request) time.Time {
	reqStartTimeRaw := r.Context().Value(consts.ReqContextStartTime)

	startTime, _ := reqStartTimeRaw.(time.Time)

	return startTime
}
