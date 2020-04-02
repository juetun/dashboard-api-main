package models

import (
	"github.com/juetun/app-dashboard/lib/base"
)

type ZPostViews struct {
	base.Model
	PostId string `gorm:"column:post_id;" json:"post_id"`
	Num    int `gorm:"column:num;" json:"num"`
}

func (r *ZPostViews) TableName() string {
	return "z_post_views"
}
