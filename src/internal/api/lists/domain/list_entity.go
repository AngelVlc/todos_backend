package domain

type ListEntity struct {
	ID         int32               `json:"id"`
	Name       ListNameValueObject `json:"name"`
	UserID     int32               `json:"-"`
	ItemsCount int32               `json:"itemsCount"`
	Items      []*ListItemEntity   `json:"items,omitempty"`
}

func (e *ListEntity) ToListRecord() *ListRecord {
	r := &ListRecord{
		ID:         e.ID,
		Name:       e.Name.String(),
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

func (e *ListEntity) GetMaxItemPosition() int32 {
	var max int32 = 0
	for _, v := range e.Items {
		if v.Position > max {
			max = v.Position
		}
	}

	return max
}
