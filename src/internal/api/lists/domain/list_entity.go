package domain

import "fmt"

type ListEntity struct {
	ID         int32               `json:"id"`
	Name       ListNameValueObject `json:"name"`
	UserID     int32               `json:"-"`
	CategoryID *int32              `json:"categoryId"`
	ItemsCount int32               `json:"itemsCount"`
	Items      []*ListItemEntity   `json:"items,omitempty"`
}

func (e *ListEntity) ToListRecord() *ListRecord {
	r := &ListRecord{
		ID:         e.ID,
		Name:       e.Name.String(),
		CategoryID: e.CategoryID,
		UserID:     e.UserID,
		ItemsCount: e.ItemsCount,
		Items:      make([]*ListItemRecord, len(e.Items)),
	}

	for i, v := range e.Items {
		r.Items[i] = &ListItemRecord{
			ID:          v.ID,
			ListID:      v.ListID,
			UserID:      v.UserID,
			Title:       v.Title.String(),
			Description: v.Description.String(),
			Position:    int32(i),
		}
	}

	return r
}

func (e *ListEntity) ToListSearchDocument() ListSearchDocument {
	d := ListSearchDocument{
		ObjectID:          fmt.Sprint(e.ID),
		UserID:            e.UserID,
		Name:              e.Name,
		ItemsTitles:       make([]string, len(e.Items)),
		ItemsDescriptions: make([]string, len(e.Items)),
	}

	for i, v := range e.Items {
		d.ItemsTitles[i] = v.Title.String()
		d.ItemsDescriptions[i] = v.Description.String()
	}

	return d
}

func (e *ListEntity) GetMaxItemPosition() int32 {
	var max int32 = 0
	for _, v := range e.Items {
		if v.Position > max {
			max = v.Position
		}
	}

	return max
}
