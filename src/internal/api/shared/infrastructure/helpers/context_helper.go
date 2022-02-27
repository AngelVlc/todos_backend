package helpers

import (
	"net/http"
	"time"

	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/consts"
)

func GetUserIDFromContext(r *http.Request) int32 {
	userIDRaw := r.Context().Value(consts.ReqContextUserIDKey)

	userID, _ := userIDRaw.(int32)

	return userID
}

func GetRequestIDFromContext(r *http.Request) string {
	requestIDRaw := r.Context().Value(consts.ReqContextRequestKey)

	requestID, _ := requestIDRaw.(string)

	return requestID
}

func GetRequestStartTimeFromContext(r *http.Request) time.Time {
	requestIDRaw := r.Context().Value(consts.ReqContextStartTime)

	startTime, _ := requestIDRaw.(time.Time)

	return startTime
}
