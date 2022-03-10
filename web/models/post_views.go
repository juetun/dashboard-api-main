package models

import (
	"fmt"
	"github.com/juetun/base-wrapper/lib/base"
)

type ZPostViews struct {
	base.Model
	PostId string `gorm:"column:post_id;" json:"post_id"`
	Num    int `gorm:"column:num;" json:"num"`
}

func (r *ZPostViews) TableName() string {
 	return fmt.Sprintf("%spost_views", TablePrefix)
}
