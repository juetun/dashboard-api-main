/**
* @Author:changjiang
* @Description:
* @File:admin_app
* @Version: 1.0.0
* @Date 2021/5/22 11:23 下午
 */
package models

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/juetun/base-wrapper/lib/base"
)

type AdminApp struct {
	Id         int               `json:"id" gorm:"column:id;primary_key" `
	UniqueKey  string            `json:"unique_key" gorm:"column:unique_key"`
	Hosts      string            `json:"-" gorm:"column:hosts""`
	HostConfig map[string]string `json:"hosts" gorm:"-"`
	Port       int               `json:"port"  gorm:"column:port"`
	Name       string            `json:"name" gorm:"column:name"`
	Desc       string            `json:"desc" gorm:"column:desc"`
	IsStop     int               `json:"is_stop" gorm:"column:is_stop"`
	CreatedAt  base.TimeNormal   `json:"created_at" gorm:"column:created_at" `
	UpdatedAt  base.TimeNormal   `json:"updated_at" gorm:"column:updated_at" `
	DeletedAt  *time.Time        `json:"-" gorm:"column:deleted_at" `
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
	return "admin_app"
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
