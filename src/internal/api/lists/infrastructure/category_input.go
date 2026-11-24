package infrastructure

import (
	"encoding/json"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
)

type CategoryInput struct {
	Name        domain.CategoryNameValueObject        `json:"name"`
	Description domain.CategoryDescriptionValueObject `json:"description"`
	IsFavourite bool                                  `json:"isFavourite"`
}

func (i *CategoryInput) UnmarshalJSON(data []byte) error {
	var realInput struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		IsFavourite bool   `json:"isFavourite"`
	}

	if err := json.Unmarshal(data, &realInput); err != nil {
		return err
	}

	nvo, err := domain.NewCategoryNameValueObject(realInput.Name)
	if err != nil {
		return err
	}

	dvo, err := domain.NewCategoryDescriptionValueObject(realInput.Description)
	if err != nil {
		return err
	}

	*i = CategoryInput{
		Name:        nvo,
		Description: dvo,
		IsFavourite: realInput.IsFavourite,
	}

	return nil
}

func (i *CategoryInput) ToCategoryEntity() *domain.CategoryEntity {
	list := &domain.CategoryEntity{
		Name:        i.Name,
		Description: i.Description,
		IsFavourite: i.IsFavourite,
	}

	return list
}
