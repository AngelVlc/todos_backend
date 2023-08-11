package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/events"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/infrastructure/results"
)

func IndexAllListsHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	go h.EventBus.Publish(events.IndexAllListsRequested, nil)

	return results.OkResult{Content: nil, StatusCode: http.StatusNoContent}
}
