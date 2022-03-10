package models

import (
	"fmt"
	"github.com/juetun/base-wrapper/lib/base"
)

type ZPostTag struct {
	base.Model
	PostId int `gorm:"column:post_id;" json:"post_id"`
	TagId  int `gorm:"column:tag_id;" json:"tag_id"`
}

func (r *ZPostTag) TableName() string {
	return fmt.Sprintf("%spost_tag", TablePrefix)
}
