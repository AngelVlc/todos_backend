package infrastructure

import "github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"

type ListInput struct {
	Name  domain.ListNameValueObject `json:"name"`
	Items []ListItemInput            `json:"items"`
}

func (i *ListInput) ToListRecord() *domain.ListRecord {
	r := &domain.ListRecord{
		Name:  i.Name.String(),
		Items: make([]*domain.ListItemRecord, len(i.Items)),
	}

	for i, v := range i.Items {
		r.Items[i] = &domain.ListItemRecord{
			ID:          v.ID,
			Title:       v.Title.String(),
			Description: v.Description.String(),
			Position:    int32(i),
		}
	}

	return r
}
