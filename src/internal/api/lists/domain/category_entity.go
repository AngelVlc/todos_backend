package domain

type CategoryEntity struct {
	ID          int32                          `json:"id"`
	Name        CategoryNameValueObject        `json:"name"`
	Description CategoryDescriptionValueObject `json:"description"`
}

func (e *CategoryEntity) ToCategoryRecord() *CategoryRecord {
	return &CategoryRecord{
		ID:          e.ID,
		Name:        e.Name.String(),
		Description: e.Description.String(),
	}
}
