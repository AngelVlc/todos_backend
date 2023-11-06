package domain

type CategoryEntity struct {
	ID          int32
	Name        CategoryNameValueObject
	Description CategoryDescriptionValueObject
}

func (e *CategoryEntity) ToCategoryRecord() *CategoryRecord {
	return &CategoryRecord{
		ID:          e.ID,
		Name:        e.Name.String(),
		Description: e.Description.String(),
	}
}
