package models

import (
	"fmt"
	"github.com/juetun/base-wrapper/lib/base"
)

type ZTags struct {
	base.Model
	Name        string `gorm:"column:name;" json:"name"`
	DisplayName string `gorm:"column:display_name;" json:"display_name"`
	SeoDesc     string `gorm:"column:seo_desc;" json:"seo_desc"`
	Num         int    `gorm:"column:num;" json:"num"`
}

func (r *ZTags) TableName() string {
	return fmt.Sprintf("%stags", TablePrefix)
}
