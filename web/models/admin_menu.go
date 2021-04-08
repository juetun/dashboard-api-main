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
	"github.com/juetun/base-wrapper/lib/hashid"
)

type AdminMenu struct {
	Id         int       `gorm:"primary_key" json:"id" form:"id"`
	PermitKey  string    `json:"permit_key" gorm:"permit_key"`
	ParentId   int       `json:"parent_id" gorm:"parent_id" form:"parent_id"`
	Module     string    `json:"module" gorm:"module" form:"module"`
	Label      string    `json:"label" gorm:"label" form:"label"`
	Icon       string    `json:"icon" gorm:"icon" form:"icon"`
	HideInMenu uint8     `json:"hide_in_menu" gorm:"hide_in_menu" form:"hide_in_menu"`
	UrlPath    string    `json:"url_path" gorm:"url_path" form:"url_path"`
	SortValue  int       `json:"sort_value" gorm:"sort_value" form:"sort_value"`
	OtherValue string    `json:"other_value" gorm :"other_value" form:"other_value"`
	CreatedAt  time.Time `json:"-" gorm :"created_at"  form:"-"`
	UpdatedAt  time.Time `json:"-" gorm :"updated_at"  form:"-"`
	IsDel      int       `json:"-" gorm:"is_del" form:"is_del"`
}

func (r *AdminMenu) TableName() string {
	return "admin_menu"
}
func (r *AdminMenu) AfterCreate(tx *gorm.DB) (err error) {
	if r.PermitKey == "" {
		tx.Model(r).Where("id=?", r.Id).Update("permit_key", r.getPathName())
	}
	return
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
