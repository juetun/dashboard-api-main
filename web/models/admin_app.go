// Package models /**
package models

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type AdminApp struct {
	Id        int        `json:"id" gorm:"column:id;primary_key" `
	UniqueKey string     `json:"unique_key" gorm:"column:unique_key"`
	Hosts     string     `json:"-" gorm:"column:hosts"`
	Port      int        `json:"port"  gorm:"column:port"`
	Name      string     `json:"name" gorm:"column:name"`
	Desc      string     `json:"desc" gorm:"column:desc"`
	IsStop    int        `json:"is_stop" gorm:"column:is_stop"`
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
