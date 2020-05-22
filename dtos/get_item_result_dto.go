package dtos

// GetItemResultDto is the model for a single list item
type GetItemResultDto struct {
	ID          int32  `json:"id" gorm:"type:int(32);primary_key"`
	Title       string `json:"title" gorm:"type:varchar(50)"`
	Description string `json:"description" gorm:"type:varchar(200)"`
}

func (GetItemResultDto) TableName() string {
	return "listItems"
}
