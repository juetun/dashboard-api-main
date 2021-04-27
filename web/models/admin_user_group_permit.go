/**
* @Author:changjiang
* @Description:
* @File:admin_user_group_permit
* @Version: 1.0.0
* @Date 2020/9/16 10:38 下午
 */
package models

import (
	"time"
)

const (
	PathTypePage = "page"
	PathTypeApi  = "api"
	SetPermitAdd ="add"
	SetPermitCancel ="cancel"
)

type AdminUserGroupPermit struct {
	Id        int        `gorm:"column:id;primary_key" json:"id"`
	GroupId   int        `json:"group_id" gorm:"column:group_id" `
	PathType  string     `json:"path_type" gorm:"column:path_type"`
	MenuId    int        `json:"menu_id" gorm:"column:menu_id"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"-" `
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"-" `
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"-"`
}

func (r *AdminUserGroupPermit) TableName() string {
	return "admin_user_group_permit"
}
