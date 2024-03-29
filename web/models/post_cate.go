package models

import (
	"fmt"
	"github.com/juetun/base-wrapper/lib/base"
)

type ZPostCate struct {
	base.Model
	PostId string `gorm:"column:post_id;" json:"post_id"`
	CateId string `gorm:"column:cate_id;" json:"cate_id"`
}

func (r *ZPostCate) TableName() string {
	return fmt.Sprintf("%spost_cate", TablePrefix)
}
