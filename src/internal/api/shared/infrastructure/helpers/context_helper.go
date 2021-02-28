package helpers

import (
	"net/http"

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
