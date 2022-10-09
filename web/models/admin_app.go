// Package models /**
package models

import (
	"encoding/json"
	"fmt"
	"github.com/juetun/base-wrapper/lib/base"
	"time"

	"gorm.io/gorm"
)

const (
	AdminAppIsStopYes uint8 = iota + 1
	AdminAppIsStopNo
)

var (
	SliceAdminAppIsStop = base.ModelItemOptions{
		{
			Label: "已启用",
			Value: AdminAppIsStopYes,
		},
		{
			Label: "已禁用",
			Value: AdminAppIsStopNo,
		},
	}
)

type AdminApp struct {
	Id        int        `json:"id" gorm:"column:id;primary_key" `
	UniqueKey string     `json:"unique_key" gorm:"column:unique_key;type:varchar(150);not null;default:'';comment:唯一的KEY"`
	Hosts     string     `json:"-" gorm:"column:hosts;type:varchar(150);not null;default:'';comment:访问地址"`
	Port      int        `json:"port"  gorm:"column:port;type:int(7);not null;default:0;"`
	Name      string     `json:"name" gorm:"column:name;type:varchar(150);not null;default:'';comment:应用名"`
	Desc      string     `json:"desc" gorm:"column:desc;type:varchar(150);not null;default:'';comment:应用描述"`
	IsStop    uint8      `json:"is_stop"  gorm:"column:is_stop;type:tinyint(2);not null;default:0;comment:是否禁用1-启用 2-禁用"`
	CreatedAt time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"-" `
	UpdatedAt time.Time  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"-" `
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"-"`

	HostConfig map[string]string `json:"hosts" gorm:"-"`
}


func (r *AdminApp) GetTableComment() (res string) {
	return "服务管理"
}

func (r *AdminApp) UnmarshalHosts() (err error) {
	r.HostConfig = make(map[string]string)
	if r.Hosts != "" {
		err = json.Unmarshal([]byte(r.Hosts), &r.HostConfig)
		return
	}
	return
}

func (r *AdminApp) TableName() string {
	return fmt.Sprintf("%sapp", TablePrefix)
}

func (r *AdminApp) MarshalHosts() (err error) {
	var bt []byte
	bt, err = json.Marshal(r.HostConfig)
	if err != nil {
		return
	}
	r.Hosts = string(bt)
	return
}
func (r *AdminApp) AfterUpdate(tx *gorm.DB) (err error) {
	// if r.GroupCode == "" {
	// 	tx.Table(r.TableName()).
	// 		Where("id=?", r.ParentId).
	// 		Update("last_child_group_code", groupCode)
	//
	// }
	return
}
func (r *AdminApp) AfterCreate(tx *gorm.DB) (err error) {
	// if r.GroupCode == "" {
	// 	groupCode := r.getGroupCode(tx)
	// 	tx.Table(r.TableName()).
	// 		Where("id=?", r.ParentId).
	// 		Update("last_child_group_code", groupCode)
	// 	tx.Model(r).Where("id=?", r.Id).Update("group_code", groupCode)
	// }
	return
}

func (r *AdminApp) ParseIsStop() (res string) {
	mapIsStop, _ := SliceAdminAppIsStop.GetMapAsKeyUint8()
	if dt, ok := mapIsStop[r.IsStop]; ok {
		res = dt
		return
	}
	res = fmt.Sprintf("未知类型(%d)", r.IsStop)
	return
}
