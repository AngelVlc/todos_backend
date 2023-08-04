package infrastructure

import "github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"

type ListInput struct {
	Name  domain.ListNameValueObject `json:"name"`
	Items []ListItemInput            `json:"items"`
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
