// Package models /**
package models

import (
	"fmt"
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

type (
	AdminImport struct {
		Id            int64      `gorm:"primary_key" json:"id" form:"id"`
		PermitKey     string     `gorm:"column:permit_key;not null;type:varchar(100);default:'';comment:KEY" json:"permit_key" form:"permit_key"`
		AppName       string     `gorm:"column:app_name;not null;type:varchar(100);default:'';comment:应用名" json:"app_name" form:"app_name"`
		AppVersion    string     `gorm:"column:app_version;not null;type:varchar(100);default:'';comment:版本" json:"app_version" form:"app_version"`
		UrlPath       string     `gorm:"column:url_path;not null;type:varchar(300);default:'';comment:接口地址" json:"url_path" form:"url_path"`
		RegexpString  string     `gorm:"column:regexp_string;not null;type:varchar(300);default:'';comment:通配地址" json:"regexp_string" form:"regexp_string"`
		SortValue     int        `gorm:"column:sort_value;not null;type:int(10);default:0;comment:排序值" json:"sort_value" form:"sort_value"`
		RequestMethod string     `gorm:"column:request_method;not null;type:varchar(100);default:'';comment:接口支持方法" json:"request_method" form:"request_method"`
		DefaultOpen   uint8      `gorm:"column:default_open;not null;type:tinyint(2);default:1;comment:默认是否开启跟随父界面设置权限1-开启 2-不开启" json:"default_open" form:"default_open"`
		NeedLogin     uint8      `gorm:"column:need_login;not null;type:tinyint(2);default:1;comment:是否需要登录1-需要 2-不需要" json:"need_login" form:"need_login"`
		NeedSign      uint8      `gorm:"column:need_sign;not null;type:tinyint(2);default:1;comment:接口是否需要签名验证1-需要 2-不需要" json:"need_sign" form:"need_sign"`
		CreatedAt     time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"-" `
		UpdatedAt     time.Time  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"-" `
		DeletedAt     *time.Time `gorm:"column:deleted_at" json:"-"`
	}

)

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
	return fmt.Sprintf("%simport", TablePrefix)
}

func (r *AdminImport) GetPathName() (res string) {
	res, _ = hashid.Encode(r.TableName(), int(r.Id))
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

// MatchPath 匹配一个方法链接是否有权限
func (r *AdminImport) MatchPath(uri string, method string) (res bool) {

	methods := r.GetRequestMethods()
	var eqMethod bool
	for _, mth := range methods {
		if mth == method {
			eqMethod = true
			break
		}
	}

	if !eqMethod { // 如果不支持方法
		return
	}

	var (
		regexp string
		err    error
	)

	if r.UrlPath == "" {
		return
	}

	if regexp, err = route_match.RoutePathToRegexp(r.UrlPath); err != nil {
		return
	}

	if res, err = route_match.RoutePathMath(regexp, uri); err != nil {
		return
	}
	return
}
