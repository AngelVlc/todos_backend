package domain

type CategoryEntity struct {
	ID          int32                          `json:"id"`
	Name        CategoryNameValueObject        `json:"name"`
	UserID      int32                          `json:"-"`
	Description CategoryDescriptionValueObject `json:"description"`
}

func (e *CategoryEntity) ToCategoryRecord() *CategoryRecord {
	return &CategoryRecord{
		ID:          e.ID,
		Name:        e.Name.String(),
		UserID:      e.UserID,
		Description: e.Description.String(),
	}
}
