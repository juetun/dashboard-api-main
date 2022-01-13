// Package models
package models

import (
	"reflect"
	"strings"
	"time"

	"github.com/juetun/base-wrapper/lib/utils/hashid"
	"gorm.io/gorm"
)

//是否显示或隐藏菜单
const (
	AdminMenuHideInMenuTrue = iota + 1
	AdminMenuHideInMenuNo
)

const (
	AdminMenuIsHomePageYes = iota + 1 //首页
	AdminMenuIsHomePageNo //不是首页
)

//是否显示或隐藏菜单
const (
	AdminMenuManageImportPermitTrue = iota + 1 //管理
	AdminMenuManageImportPermitNo              //不管理
)

const (
	// AdminMenuNameMaxLength 菜单名称最大长度
	AdminMenuNameMaxLength = 6

	CommonMenuDefaultHomePage = "首页"
	CommonMenuDefaultLabel    = "公共接口" // 每个子系统的公共接口菜单名称（开启了系统首页就会有此菜单功能）
)

type AdminMenu struct {
	Id                 int64      `gorm:"column:id;primary_key" json:"id" form:"id"`
	PermitKey          string     `json:"permit_key" gorm:"column:permit_key" form:"permit_key"`
	ParentId           int64      `json:"parent_id" gorm:"column:parent_id" form:"parent_id"`
	Module             string     `json:"module" gorm:"column:module" form:"module"`
	Label              string     `json:"label" gorm:"column:label" form:"label"`
	Icon               string     `json:"icon" gorm:"column:icon" form:"icon"`
	IsHomePage         uint8      `gorm:"column:is_home_page;default:0;not null;type:tinyint(2);default:2;comment:是否为首页1-不是,2-是" json:"is_home_page" form:"is_home_page"`
	HideInMenu         uint8      `gorm:"column:hide_in_menu;default:0;not null;type:tinyint(2);default:1;comment:菜单是否显示 1-显示,2-不显示" json:"hide_in_menu" form:"hide_in_menu"`
 	ManageImportPermit uint8      `gorm:"column:manage_import_permit;default:1;not null;type:tinyint(2);default:0;comment:是否有管理接口的权限1-管理,2-不管理" json:"manage_import_permit" form:"manage_import_permit"`
	Domain             string     `json:"domain" gorm:"column:domain" form:"domain"`
	UrlPath            string     `json:"url_path" gorm:"column:url_path" form:"url_path"`
	SortValue          int        `json:"sort_value" gorm:"column:sort_value" form:"sort_value"`
	OtherValue         string     `json:"other_value" gorm:"column:other_value" form:"other_value"`
	CreatedAt          time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"-" `
	UpdatedAt          time.Time  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"-" `
	DeletedAt          *time.Time `gorm:"column:deleted_at" json:"-"`
}

func (r *AdminMenu) TableName() string {
	return "admin_menu"
}

func (r *AdminMenu) getPathName() (res string) {
	res, _ = hashid.Encode(r.TableName(), int(r.Id))
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

type DefaultSystemMenuNeedParams struct {
	Module         string    `json:"module"`
	UpdateTime     time.Time `json:"update_time"`
	CreateTime     time.Time `json:"create_time"`
	ParentSystemId int64     `json:"parent_system_id"`
}

func (r *AdminMenu) InitDefaultSystemMenu(data *DefaultSystemMenuNeedParams) (res []*AdminMenu) {

	res = []*AdminMenu{
		{
			Module:             data.Module,
			ParentId:           data.ParentSystemId,
			Label:              CommonMenuDefaultHomePage,
			Icon:               "md-home",
			ManageImportPermit: AdminMenuManageImportPermitTrue,
			HideInMenu:         AdminMenuHideInMenuNo,
			IsHomePage:         AdminMenuIsHomePageYes,
			UrlPath:            "",
			SortValue:          90000001,
			OtherValue:         "",
			CreatedAt:          data.CreateTime,
			UpdatedAt:          data.UpdateTime,
		},
		{
			Module:             data.Module,
			ParentId:           data.ParentSystemId,
			Label:              CommonMenuDefaultLabel,
			Icon:               "",
			ManageImportPermit: 1,
			HideInMenu:         1,
			UrlPath:            "",
			SortValue:          90000000,
			OtherValue:         "",
			CreatedAt:          data.CreateTime,
			UpdatedAt:          data.UpdateTime,
		},
	}
	return
}
