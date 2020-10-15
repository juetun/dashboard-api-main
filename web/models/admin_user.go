/**
* @Author:changjiang
* @Description:
* @File:admin_user
* @Version: 1.0.0
* @Date 2020/9/16 10:38 下午
 */
package models

import (
	"github.com/juetun/base-wrapper/lib/base"
)

type AdminUser struct {
	Id        int             `gorm:"primary_key" json:"id"`
	UserHid   string          `json:"user_hid" gorm:"user_hid"`
	RealName  string          `json:"real_name" gorm:"real_name"`
	CreatedAt base.TimeNormal `json:"created_at" gorm:"created_at" `
	UpdatedAt base.TimeNormal `json:"updated_at" gorm:"updated_at" `
	DeletedAt base.TimeNormal `json:"-" gorm:"deleted_at" `
}

func (r *AdminUser) TableName() string {
	return "admin_user"
}
