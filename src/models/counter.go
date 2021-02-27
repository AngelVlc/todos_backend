package models

type Counter struct {
	ID    int32  `gorm:"type:int(32);primary_key"`
	Name  string `gorm:"type:varchar(10);index:idx_counters_name"`
	Value int32  `gorm:"type:int(32)"`
}
