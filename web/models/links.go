package models

import (
	"github.com/juetun/base-wrapper/lib/base"
)

type ZLinks struct {
	base.Model
	Name  string `gorm:"column:name;" json:"name"`
	Link  string `gorm:"column:link;" json:"link"`
	Order int    `gorm:"column:order;" json:"order"`
}

func (r *ZLinks) TableName() string {
	return "z_links"
}
