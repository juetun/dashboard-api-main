/**
* @Author:changjiang
* @Description:
* @File:admin_app
* @Version: 1.0.0
* @Date 2021/5/22 11:23 下午
 */
package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/juetun/base-wrapper/lib/base"
)

type AdminApp struct {
	Id        int    `json:"id" gorm:"column:id;primary_key" `
	UniqueKey string `json:"unique_key" gorm:"column:unique_key"`
	Port      int    `json:"port"  gorm:"column:port"`
	Name      string `json:"name" gorm:"column:name"`
	Desc      string `json:"desc" gorm:"column:desc"`
	IsStop    uint8  `json:"is_stop" gorm:"column:is_stop"`
	CreatedAt base.TimeNormal `json:"created_at" gorm:"column:created_at" `
	UpdatedAt base.TimeNormal `json:"updated_at" gorm:"column:updated_at" `
	DeletedAt *time.Time      `json:"-" gorm:"column:deleted_at" `
}

func (r *AdminApp) TableName() string {
	return "admin_group"
}

func (r AdminApp) AfterUpdate(tx *gorm.DB) (err error) {
	// if r.GroupCode == "" {
	// 	tx.Table(r.TableName()).
	// 		Where("id=?", r.ParentId).
	// 		Update("last_child_group_code", groupCode)
	//
	// }
	return
}
func (r AdminApp) AfterCreate(tx *gorm.DB) (err error) {
	// if r.GroupCode == "" {
	// 	groupCode := r.getGroupCode(tx)
	// 	tx.Table(r.TableName()).
	// 		Where("id=?", r.ParentId).
	// 		Update("last_child_group_code", groupCode)
	// 	tx.Model(r).Where("id=?", r.Id).Update("group_code", groupCode)
	// }
	return
}
