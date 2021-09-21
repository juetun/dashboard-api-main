/**
* @Author:changjiang
* @Description:
* @File:admin_import
* @Version: 1.0.0
* @Date 2020/12/9 11:08 下午
 */
package models

import (
	"strings"
	"time"

	"github.com/juetun/base-wrapper/lib/utils/hashid"
	"gorm.io/gorm"
)

const (
	DefaultOpen = 1
)

type AdminImport struct {
	Id            int        `gorm:"primary_key" json:"id" form:"id"`
	PermitKey     string     `json:"permit_key" gorm:"column:permit_key" form:"permit_key"`
	AppName       string     `json:"app_name" gorm:"column:app_name" form:"app_name"`
	AppVersion    string     `json:"app_version" gorm:"column:app_version" form:"app_version"`
	UrlPath       string     `json:"url_path" gorm:"column:url_path" form:"url_path"`
	SortValue     int        `json:"sort_value" gorm:"column:sort_value" form:"sort_value"`
	RequestMethod string     `json:"request_method" gorm:"column:request_method" form:"request_method"`
	DefaultOpen   uint8      `json:"default_open" gorm:"column:default_open" form:"default_open"`
	NeedLogin     uint8      `json:"need_login" gorm:"column:need_login" form:"need_login"`
	NeedSign      uint8      `json:"need_sign" gorm:"column:need_sign" form:"need_sign"`
	CreatedAt     time.Time  `json:"-" gorm:"column:created_at" `
	UpdatedAt     time.Time  `json:"-" gorm:"column:updated_at" `
	DeletedAt     *time.Time `json:"-" gorm:"column:deleted_at"`
}

func (r *AdminImport) TableName() string {
	return "admin_import"
}
func (r *AdminImport) GetPathName() (res string) {
	res, _ = hashid.Encode(r.TableName(), r.Id)
	return
}
func (r *AdminImport) AfterUpdate(tx *gorm.DB) (err error) {
	if r.PermitKey == "" {
		tx.Table(r.TableName()).
			Where("id=?", r.Id).
			Update("permit_key", r.GetPathName())
	}
	return
}
func (r *AdminImport) AfterCreate(tx *gorm.DB) (err error) {
	if r.PermitKey == "" {
		tx.Model(r).Where("id=?", r.Id).Update("permit_key", r.GetPathName())
	}
	return
}


func (r *AdminImport) SetRequestMethods(methods []string) {
	if len(methods) > 0 {
		r.RequestMethod = strings.Join(methods, ",")
	}
	return
}
func (r *AdminImport) GetRequestMethods() (res []string) {
	res = []string{}
	if r.RequestMethod != "" {
		res = strings.Split(r.RequestMethod, ",")
	}
	return
}
