/**
* @Author:changjiang
* @Description:
* @File:admin_user_group
* @Version: 1.0.0
* @Date 2020/9/16 10:38 下午
 */
package models

import (
	"time"

	"github.com/juetun/base-wrapper/lib/base"
)

type AdminUserGroup struct {
	Id        int             `gorm:"column:id;primary_key" json:"id"`
	GroupId   int             `json:"group_id" gorm:"column:group_id"`
	UserHid   string          `json:"user_hid"  gorm:"column:user_hid"`
	CreatedAt base.TimeNormal `json:"-" gorm:"column:created_at" `
	UpdatedAt base.TimeNormal `json:"updated_at" gorm:"column:updated_at" `
	DeletedAt *time.Time      `json:"-" gorm:"column:deleted_at" `
}

func (r *AdminUserGroup) TableName() string {
	return "admin_user_group"
}
