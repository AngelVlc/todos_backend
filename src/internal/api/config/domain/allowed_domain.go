package domain

type AllowedOrigin struct {
	ID     int32  `gorm:"type:int(32);primary_key"`
	Origin Origin `gorm:"column:origin;type:varchar(250)"`
}

func (AllowedOrigin) TableName() string {
	return "allowed_origins"
}
