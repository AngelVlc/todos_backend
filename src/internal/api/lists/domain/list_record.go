package domain

import "database/sql"

type ListRecord struct {
	ID         int32             `gorm:"type:int(32);primary_key"`
	Name       string            `gorm:"type:varchar(50)"`
	UserID     int32             `gorm:"column:userId;type:int(32)"`
	CategoryID *sql.NullInt32    `gorm:"column:categoryId;type:int(32)"`
	ItemsCount int32             `gorm:"column:itemsCount;type:int(32)"`
	Items      []*ListItemRecord `gorm:"foreignKey:ListID"`
}

func (ListRecord) TableName() string {
	return "lists"
}

func (r *ListRecord) ToListEntity() *ListEntity {
	nvo, _ := NewListNameValueObject(r.Name)

	items := make([]*ListItemEntity, len(r.Items))

	for i, v := range r.Items {
		tvo, _ := NewItemTitleValueObject(v.Title)
		dvo, _ := NewItemDescriptionValueObject(v.Description)

		items[i] = &ListItemEntity{
			ID:          v.ID,
			ListID:      v.ListID,
			UserID:      v.UserID,
			Title:       tvo,
			Description: dvo,
			Position:    v.Position,
		}
	}

	var categoryID int32

	if r.CategoryID != nil && r.CategoryID.Valid {
		categoryID = r.CategoryID.Int32
	}

	return &ListEntity{
		ID:         r.ID,
		Name:       nvo,
		CategoryID: &categoryID,
		UserID:     r.UserID,
		ItemsCount: r.ItemsCount,
		Items:      items,
	}
}
