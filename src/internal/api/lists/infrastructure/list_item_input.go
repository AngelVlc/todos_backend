package infrastructure

import (
	"encoding/json"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
)

type ListItemInput struct {
	ID          int32                             `json:"id"`
	Title       domain.ItemTitleValueObject       `json:"title"`
	Description domain.ItemDescriptionValueObject `json:"description"`
	Position    int32                             `json:"position"`
}

func (i *ListItemInput) UnmarshalJSON(data []byte) error {
	var realInput struct {
		ID          int32  `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Position    int32  `json:"position"`
	}

	if err := json.Unmarshal(data, &realInput); err != nil {
		return err
	}

	tvo, err := domain.NewItemTitleValueObject(realInput.Title)
	if err != nil {
		return err
	}

	dvo, err := domain.NewItemDescriptionValueObject(realInput.Description)
	if err != nil {
		return err
	}

	*i = ListItemInput{
		ID:          realInput.ID,
		Title:       tvo,
		Description: dvo,
		Position:    realInput.Position,
	}

	return nil
}
