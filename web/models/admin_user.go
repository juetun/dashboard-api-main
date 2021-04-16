/**
* @Author:changjiang
* @Description:
* @File:admin_user
* @Version: 1.0.0
* @Date 2020/9/16 10:38 下午
 */
package models

import (
	"github.com/jinzhu/gorm"
)

type AdminUser struct {
	// Id        int             `gorm:"column:id;primary_key" json:"id"`
	UserHid   string `json:"user_hid" gorm:"column:user_hid"`
	RealName  string `json:"real_name" gorm:"column:real_name"`
	Mobile    string `json:"mobile" gorm:"column:mobile"`
	FlagAdmin int    `json:"flag_admin" gorm:"column:flag_admin"`
	gorm.Model
	// CreatedAt base.TimeNormal `json:"created_at" gorm:"column:created_at" `
	// UpdatedAt base.TimeNormal `json:"updated_at" gorm:"column:updated_at" `
	// DeletedAt *time.Time      `json:"-" gorm:"column:deleted_at" `
}

func (r *AdminUser) TableName() string {
	return "admin_user"
}
