package dtos

// GetListsResultDto is the struct used as result for the Get lists method
type GetListsResultDto struct {
	ID   string `json:"id" gorm:"type:varchar(10)"`
	Name string `json:"name" gorm:"foreignkey:ListID"`
}

func (GetListsResultDto) TableName() string {
	return "lists"
}
