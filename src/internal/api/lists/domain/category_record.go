package domain

type CategoryRecord struct {
	ID          int32  `gorm:"type:int(32);primary_key"`
	Name        string `gorm:"type:varchar(12)"`
	Description string `gorm:"type:varchar(200)"`
	UserID      int32  `gorm:"column:userId;type:int(32)"`
}

func (CategoryRecord) TableName() string {
	return "categories"
}

func (r *CategoryRecord) ToCategoryEntity() *CategoryEntity {
	nvo, _ := NewCategoryNameValueObject(r.Name)
	dvo, _ := NewCategoryDescriptionValueObject(r.Description)

	return &CategoryEntity{
		ID:          r.ID,
		Name:        nvo,
		Description: dvo,
		UserID:      r.UserID,
	}
}
