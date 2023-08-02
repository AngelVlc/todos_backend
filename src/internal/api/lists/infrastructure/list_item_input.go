package infrastructure

import "github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"

type ListItemInput struct {
	ID          int32                             `json:"id"`
	Title       domain.ItemTitleValueObject       `json:"title"`
	Description domain.ItemDescriptionValueObject `json:"description"`
	Position    int32                             `json:"position"`
}
