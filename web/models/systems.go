package models

import (
	"fmt"
	"github.com/juetun/base-wrapper/lib/base"
)

type ZBaseSys struct {
	base.Model
	Theme        int    `gorm:"column:theme;" json:"theme"`
	Title        string `gorm:"column:title;" json:"title"`
	Keywords     string `gorm:"column:keywords;" json:"keywords"`
	Description  string `gorm:"column:description;" json:"description"`
	RecordNumber string `gorm:"column:record_number;" json:"record_number"`
}

func (r *ZBaseSys) TableName() string {
 	return fmt.Sprintf("%sbase_sys", TablePrefix)
}
