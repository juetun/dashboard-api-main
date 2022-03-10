package models

import (
	"fmt"
	"github.com/juetun/base-wrapper/lib/base"
)

type ZCategories struct {
	base.Model
	Name        string `gorm:"column:name;" json:"name"`
	DisplayName string `gorm:"column:display_name;" json:"display_name"`
	SeoDesc     string `gorm:"column:seo_desc;" json:"seo_desc"`
	ParentId    int    `gorm:"column:parent_id;" json:"parent_id"`
}

func (r *ZCategories) TableName() string {
	return fmt.Sprintf("%scategories", TablePrefix)
}
