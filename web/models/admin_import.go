// Package models /**
package models

import (
	"strings"
	"time"

	"github.com/juetun/base-wrapper/lib/utils/hashid"
	"github.com/juetun/dashboard-api-main/pkg/route_match"
	"gorm.io/gorm"
)

const (
	DefaultOpen = 1

	NeedLoginTrue = 1
	NeedLoginNot  = 2
	NeedSignTrue  = 1
	NeedSignNot   = 2
)

type AdminImport struct {
	Id            int        `gorm:"primary_key" json:"id" form:"id"`
	PermitKey     string     `json:"permit_key" gorm:"column:permit_key" form:"permit_key"`
	AppName       string     `json:"app_name" gorm:"column:app_name" form:"app_name"`
	AppVersion    string     `json:"app_version" gorm:"column:app_version" form:"app_version"`
	UrlPath       string     `json:"url_path" gorm:"column:url_path" form:"url_path"`
	RegexpString  string     `json:"regexp_string" gorm:"column:regexp_string" form:"regexp_string"`
	SortValue     int        `json:"sort_value" gorm:"column:sort_value" form:"sort_value"`
	RequestMethod string     `json:"request_method" gorm:"column:request_method" form:"request_method"`
	DefaultOpen   uint8      `json:"default_open" gorm:"column:default_open" form:"default_open"`
	NeedLogin     uint8      `json:"need_login" gorm:"column:need_login" form:"need_login"`
	NeedSign      uint8      `json:"need_sign" gorm:"column:need_sign" form:"need_sign"`
	CreatedAt     time.Time  `json:"-" gorm:"column:created_at" `
	UpdatedAt     time.Time  `json:"-" gorm:"column:updated_at" `
	DeletedAt     *time.Time `json:"-" gorm:"column:deleted_at"`
}

func (r *AdminImport) GetTableComment() (res string) {
	return "接口管理表"
}

func (r *AdminImport) InitRegexpString() (err error) {
	if r.UrlPath != "" {
		r.RegexpString, err = route_match.RoutePathToRegexp(r.UrlPath)
	}
	return
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
func (r *AdminImport) GetRequestMethodMap() (res map[string]uint8) {
	res = map[string]uint8{}
	if r.RequestMethod != "" {
		tmp := strings.Split(r.RequestMethod, ",")
		for _, s := range tmp {
			if s == "" {
				continue
			}
			res[strings.ToUpper(s)] = 1
		}
	}
	return
}
func (r *AdminImport) GetRequestMethods() (res []string) {
	res = []string{}
	if r.RequestMethod != "" {
		tmp := strings.Split(r.RequestMethod, ",")
		for _, s := range tmp {
			if s == "" {
				continue
			}
			res = append(res, strings.ToUpper(s))
		}
	}
	return
}
