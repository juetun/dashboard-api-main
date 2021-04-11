/**
* @Author:changjiang
* @Description:
* @File:admin_group
* @Version: 1.0.0
* @Date 2020/9/20 6:40 下午
 */
package models

import (
	"time"

	"github.com/juetun/base-wrapper/lib/base"
)

type AdminGroup struct {
	Id           int             `gorm:"primary_key" json:"id"`
	Name         string          `json:"name" gorm:"name"`
	IsSuperAdmin uint8           `json:"is_super_admin"`
	IsDel        int             `json:"-" gorm:"is_del"`
	CreatedAt    base.TimeNormal `json:"created_at" gorm:"created_at" `
	UpdatedAt    base.TimeNormal `json:"updated_at" gorm:"updated_at" `
	DeletedAt    *time.Time      `json:"-" gorm:"deleted_at" `
}

func (r *AdminGroup) TableName() string {
	return "admin_group"
}
