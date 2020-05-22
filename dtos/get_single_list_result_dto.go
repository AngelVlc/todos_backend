package dtos

// GetSingleListResultDto is the struct used as result for the Get single list method
type GetSingleListResultDto struct {
	ID        int32                        `json:"id" gorm:"type:int(32);primary_key"`
	Name      string                       `json:"name" gorm:"type:varchar(50)"`
	ListItems []GetSingleListResultItemDto `json:"items" gorm:"foreignkey:ListID"`
}

func (GetSingleListResultDto) TableName() string {
	return "lists"
}

// GetSingleListResultItemDto is the model for a single list item
type GetSingleListResultItemDto struct {
	ID          int32  `json:"id" gorm:"type:int(32);primary_key"`
	ListID      int32  `json:"listId" gorm:"column:listId;type:int(32)"`
	Title       string `json:"title" gorm:"type:varchar(50)"`
	Description string `json:"description" gorm:"type:varchar(200)"`
}

func (GetSingleListResultItemDto) TableName() string {
	return "listItems"
}
