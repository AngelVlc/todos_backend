package infrastructure

import (
	"encoding/json"
	"fmt"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
)

type ListInput struct {
	Name  domain.ListNameValueObject `json:"name"`
	Items []ListItemInput            `json:"items"`
}

func (i *ListInput) UnmarshalJSON(data []byte) error {
	var realInput struct {
		Name  string `json:"name"`
		Items []struct {
			ID          int32  `json:"id"`
			Title       string `json:"title"`
			Description string `json:"description"`
			Position    int32  `json:"position"`
		} `json:"items"`
	}

	if err := json.Unmarshal(data, &realInput); err != nil {
		return err
	}

	nvo, err := domain.NewListNameValueObject(realInput.Name)
	if err != nil {
		return err
	}

	*i = ListInput{
		Name:  nvo,
		Items: make([]ListItemInput, len(realInput.Items)),
	}

	for index, v := range realInput.Items {
		tvo, err := domain.NewItemTitleValueObject(v.Title)
		if err != nil {
			return fmt.Errorf("Item #%v: %v", index, err)
		}

		dvo, err := domain.NewItemDescriptionValueObject(v.Description)
		if err != nil {
			return fmt.Errorf("Item #%v: %v", index, err)
		}

		i.Items[index] = ListItemInput{
			ID:          v.ID,
			Title:       tvo,
			Description: dvo,
			Position:    v.Position,
		}
	}

	return nil
}

func (i *ListInput) ToListEntity() *domain.ListEntity {
	list := &domain.ListEntity{
		Name:  i.Name,
		Items: make([]*domain.ListItemEntity, len(i.Items)),
	}

	for i, v := range i.Items {
		list.Items[i] = &domain.ListItemEntity{
			ID:          v.ID,
			Title:       v.Title,
			Description: v.Description,
			Position:    int32(i),
		}
	}

	return list
}
