package models

import (
	"fmt"
	"github.com/juetun/base-wrapper/lib/base"
)

type ZPosts struct {
	base.Model
	Uid      string `gorm:"column:uid;" json:"uid"`
	UserHId  int64  `gorm:"column:user_hid;" json:"user_id"`
	Title    string `gorm:"column:title;" json:"title"`
	Summary  string `gorm:"column:summary;" json:"summary"`
	Original string `gorm:"column:original;" json:"original"`
	Content  string `gorm:"column:content;" json:"content"`
	Password string `gorm:"column:password;" json:"password"`
}

func (r *ZPosts) TableName() string {
 	return fmt.Sprintf("%sposts", TablePrefix)
}
