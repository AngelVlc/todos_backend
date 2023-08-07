package helpers

import (
	"net/http"

	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/consts"
)

func GetRequestIDFromContext(r *http.Request) string {
	requestIDRaw := r.Context().Value(consts.ReqContextRequestKey)

	requestID, _ := requestIDRaw.(string)

	return requestID
}
