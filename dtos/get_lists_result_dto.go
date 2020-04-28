package dtos

// GetListsResultDto is the struct used as result for the Get lists method
type GetListsResultDto struct {
	ID   int32  `json:"id" gorm:"type:int(32);primary_key"`
	Name string `json:"name" gorm:"type:varchar(50)"`
}

func (GetListsResultDto) TableName() string {
	return "lists"
}
