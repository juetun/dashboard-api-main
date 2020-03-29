package models

import (
	"github.com/juetun/app-dashboard/lib/base"
)

type ZPostCate struct {
	base.Model
	PostId int `gorm:"column:post_id;" json:"post_id"`
	CateId int `gorm:"column:cate_id;" json:"cate_id"`
}

func (r *ZPostCate) TableName() string {
	return "z_post_cate"
}
