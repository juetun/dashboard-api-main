package models

import (
	"github.com/juetun/app-dashboard/lib/base"
)

type ZPostTag struct {
	base.Model
	PostId int `gorm:"column:post_id;" json:"post_id"`
	TagId  int `gorm:"column:tag_id;" json:"tag_id"`
}

func (r *ZPostTag) TableName() string {
	return "z_post_tag"
}
