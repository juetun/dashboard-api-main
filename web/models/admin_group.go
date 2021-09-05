/**
* @Author:changjiang
* @Description:
* @File:admin_group
* @Version: 1.0.0
* @Date 2020/9/20 6:40 下午
 */
package models

import (
	"strings"
	"time"

	"github.com/juetun/base-wrapper/lib/base"
	"gorm.io/gorm"
)

const (
	GroupCodePrefix     = "A2" // 权限用户组统一前缀
	GroupCodeStepLength = 3    // 权限用户组统一长度


	MAXGroupNameLength  = 20   // 组名长度不能超过20个汉字
)

type AdminGroup struct {
	Id                 int             `json:"id" gorm:"column:id;primary_key" `
	Name               string          `json:"name" gorm:"column:name"`
	ParentId           int             `json:"parent_id"  gorm:"column:parent_id"`
	GroupCode          string          `json:"group_code" gorm:"column:group_code"`
	LastChildGroupCode string          `json:"last_child_group_code" gorm:"column:last_child_group_code"`
	IsSuperAdmin       uint8           `json:"is_super_admin" gorm:"column:is_super_admin"`
	IsAdminGroup       uint8           `json:"is_admin_group" gorm:"column:is_admin_group"`
	CreatedAt          base.TimeNormal `json:"created_at" gorm:"column:created_at" `
	UpdatedAt          base.TimeNormal `json:"updated_at" gorm:"column:updated_at" `
	DeletedAt          *time.Time      `json:"-" gorm:"column:deleted_at" `
}

func (r *AdminGroup) TableName() string {
	return "admin_group"
}

func (r AdminGroup) AfterUpdate(tx *gorm.DB) (err error) {
	if r.GroupCode == "" {
		groupCode := r.getGroupCode(tx)
		tx.Table(r.TableName()).
			Where("id=?", r.ParentId).
			Update("last_child_group_code", groupCode)
		tx.Table(r.TableName()).
			Where("id=?", r.Id).
			Update("group_code", groupCode)
	}
	return
}
func (r AdminGroup) AfterCreate(tx *gorm.DB) (err error) {
	if r.GroupCode == "" {
		groupCode := r.getGroupCode(tx)
		tx.Table(r.TableName()).
			Where("id=?", r.ParentId).
			Update("last_child_group_code", groupCode)
		tx.Model(r).Where("id=?", r.Id).Update("group_code", groupCode)
	}
	return
}

func (r *AdminGroup) getGroupCodePrefix(tx *gorm.DB) (prefix string) {
	prefix = GroupCodePrefix
	if r.ParentId != 0 {
		var dt []AdminGroup
		tx.Model(r).
			Where("id=?", r.ParentId).
			Order("group_code desc").
			Limit(1).
			Find(&dt)
		prefix = dt[0].GroupCode
	}
	return
}
func (r *AdminGroup) getGroupCode(tx *gorm.DB) (res string) {
	prefix := r.getGroupCodePrefix(tx)
	var fromString []byte
	var dt []AdminGroup
	tx.Model(r).
		Where("id=?", r.ParentId).
		Order("group_code desc").
		Limit(1).
		Find(&dt)
	if len(dt) > 0 {
		fromString = []byte(strings.TrimPrefix(dt[0].LastChildGroupCode, prefix))
	}
	res = prefix + r.orgString(fromString)
	return
	// GroupCodePrefix     = "A2" // 权限用户组统一前缀
	// GroupCodeStepLength = 4    // 权限用户组统一长度
}
func (r *AdminGroup) orgString(fromString []byte) (rs string) {
	s := []byte(`A0BCDE12FGHIJKLMNOP34QRST56UVW7XYZ89`)
	code := make([]byte, 0, GroupCodeStepLength)
	var s1 []byte
	fromCode := r.flipByte(fromString)
	var l = len(s)
	var nextNeedAdd bool
	var loc int
	for i := 0; i < GroupCodeStepLength; i++ {
		s1 = append(s[i+GroupCodeStepLength:], s[:i+GroupCodeStepLength]...)
		if len(fromString) == 0 {
			code = append(code, s1[0])
			continue
		}
		if i == 0 {
			if loc = r.getIndex(fromCode[i], s1) + 1; loc == l {
				nextNeedAdd = true
				loc = 0
			}
			code = append(code, s1[loc])
		} else if nextNeedAdd {
			nextNeedAdd = false
			if loc = r.getIndex(fromCode[i], s1) + 1; loc == l {
				nextNeedAdd = true
				loc = 0
			}
			code = append(code, s1[loc])
		} else {
			code = append(code, fromCode[i])
		}

	}
	rs = string(r.flipByte(code))
	return
}
func (r *AdminGroup) flipByte(arg []byte) (res []byte) {
	l := len(arg)
	res = make([]byte, l)
	for i := 0; i < l; i++ {
		res[l-i-1] = arg[i]
	}
	return
}
func (r *AdminGroup) getIndex(d byte, f []byte) (res int) {
	res = -1
	for k, value := range f {
		if d == value {
			res = k
			break
		}
	}
	return
}
