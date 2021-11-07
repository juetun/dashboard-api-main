// Package models
package models

import (
	"reflect"
	"strings"
	"time"

	"github.com/juetun/base-wrapper/lib/utils/hashid"
	"gorm.io/gorm"
)

const (
	// AdminMenuNameMaxLength 菜单名称最大长度
	AdminMenuNameMaxLength = 6
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
	Domain             string     `json:"domain" gorm:"column:domain" form:"domain"`
	UrlPath            string     `json:"url_path" gorm:"column:url_path" form:"url_path"`
	SortValue          int        `json:"sort_value" gorm:"column:sort_value" form:"sort_value"`
	OtherValue         string     `json:"other_value" gorm:"column:other_value" form:"other_value"`
	CreatedAt          time.Time  `json:"-" gorm:"column:created_at"  form:"-"`
	UpdatedAt          time.Time  `json:"-" gorm:"column:updated_at"  form:"-"`
	DeletedAt          *time.Time `json:"-" gorm:"column:deleted_at"`
}

func (r *AdminMenu) TableName() string {
	return "admin_menu"
}

func (r *AdminMenu) getPathName() (res string) {
	res, _ = hashid.Encode(r.TableName(), r.Id)
	return
}
func (r *AdminMenu) ToMapStringInterface() (res map[string]interface{}) {

	var tName string
	var tag string

	types := reflect.TypeOf(*r)
	values := reflect.ValueOf(*r)

	var fieldNum = types.NumField()
	res = make(map[string]interface{}, fieldNum)
	for i := 0; i < fieldNum; i++ {
		tName = types.Field(i).Tag.Get("gorm")
		tag = r.getColumnName(tName)
		if tag == "id" || tag == "created_at" || tag == "" {
			continue
		}
		res[tag] = values.Field(i).Interface()
	}

	return
}
func (r *AdminMenu) getColumnName(s string) (res string) {
	li := strings.Split(s, ";")
	res = s
	for _, s2 := range li {
		if s2 == "" {
			return
		}
		li1 := strings.Split(s2, ":")
		if len(li1) > 1 && li1[0] == "column" {
			res = li1[1]
		}
	}
	return
}

// AfterUpdate Updating data in same transaction
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
