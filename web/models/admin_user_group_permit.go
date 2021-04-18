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

type AdminUserGroupPermit struct {
	Id        int        `gorm:"column:id;primary_key" json:"id"`
	GroupId   int        `json:"group_id" gorm:"column:group_id" `
	PathType  string     `json:"path_type" gorm:"column:path_type"`
	MenuId    int        `json:"menu_id" gorm:"column:menu_id"`
	CreatedAt time.Time  `json:"-" gorm :"column:created_at" `
	UpdatedAt time.Time  `json:"-" gorm :"column:updated_at" `
	DeletedAt *time.Time `json:"-" gorm:"column:deleted_at"`
}

func (r *AdminUserGroupPermit) TableName() string {
	return "admin_user_group_permit"
}
