/**
* @Author:changjiang
* @Description:
* @File:admin_menu
* @Version: 1.0.0
* @Date 2020/9/16 10:37 下午
 */
package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/juetun/base-wrapper/lib/utils/hashid"
)

type AdminMenu struct {
	Id                 int        `gorm:"column:id;primary_key" json:"id" form:"id"`
	PermitKey          string     `json:"permit_key" gorm:"column:permit_key" form:"permit_key"`
	ParentId           int        `json:"parent_id" gorm:"column:parent_id" form:"parent_id"`
	Module             string     `json:"module" gorm:"column:module" form:"module"`
	Label              string     `json:"label" gorm:"column:label" form:"label"`
	Icon               string     `json:"icon" gorm:"column:icon" form:"icon"`
	HideInMenu         uint8      `json:"hide_in_menu" gorm:"column:hide_in_menu" form:"hide_in_menu"`
	ManageImportPermit uint8      `json:"manage_import_permit" gorm:"column:manage_import_permit" form:"manage_import_permit"`
	UrlPath            string     `json:"url_path" gorm:"column:url_path" form:"url_path"`
	SortValue          int        `json:"sort_value" gorm:"column:sort_value" form:"sort_value"`
	OtherValue         string     `json:"other_value" gorm :"column:other_value" form:"other_value"`
	CreatedAt          time.Time  `json:"-" gorm :"column:created_at"  form:"-"`
	UpdatedAt          time.Time  `json:"-" gorm :"column:updated_at"  form:"-"`
	DeletedAt          *time.Time `json:"-" gorm:"column:deleted_at"`
}

func (r *AdminMenu) TableName() string {
	return "admin_menu"
}

func (r *AdminMenu) getPathName() (res string) {
	res, _ = hashid.Encode(r.TableName(), r.Id)
	return
}

// Updating data in same transaction
func (r *AdminMenu) AfterUpdate(tx *gorm.DB) (err error) {
	if r.PermitKey == "" {
		tx.Table(r.TableName()).
			Where("id=?", r.Id).
			Update("permit_key", r.getPathName())
	}
	return
}
func (r *AdminMenu) AfterCreate(tx *gorm.DB) (err error) {
	if r.PermitKey == "" {
		tx.Model(r).Where("id=?", r.Id).Update("permit_key", r.getPathName())
	}
	return
}
