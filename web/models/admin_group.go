// Package models
package models

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

const (
	GroupCodePrefix        = "2" // 权限用户组统一前缀
	GroupCodeTopStepLength = 8   // 第一级用户组统一长度
	GroupCodeStepLength    = 3   // 权限用户组（非第一级）统一长度
	MAXGroupNameLength     = 80  // 组名长度不能超过20个字符
)

type AdminGroup struct {
	Id                 int64      `json:"id" gorm:"column:id;primary_key" `
	Name               string     `json:"name" gorm:"column:name"`
	ParentId           int64      `json:"parent_id"  gorm:"column:parent_id"`
	GroupCode          string     `json:"group_code" gorm:"column:group_code"`
	LastChildGroupCode string     `json:"last_child_group_code" gorm:"column:last_child_group_code"`
	IsSuperAdmin       uint8      `json:"is_super_admin" gorm:"column:is_super_admin"`
	IsAdminGroup       uint8      `json:"is_admin_group" gorm:"column:is_admin_group"`
	CreatedAt          time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"-" `
	UpdatedAt          time.Time  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"-" `
	DeletedAt          *time.Time `gorm:"column:deleted_at" json:"-"`
}

func (r *AdminGroup) GetTableComment() (res string) {
	return "管理员组表"
}

func (r *AdminGroup) TableName() string {
	return fmt.Sprintf("%sgroup", TablePrefix)
}

func (r AdminGroup) AfterUpdate(tx *gorm.DB) (err error) {
	if r.GroupCode == "" {
		var (
			groupCodeV      string
			parentGroupCode string
		)
		if groupCodeV, parentGroupCode, err = r.getGroupCode(tx); err != nil {
			return
		}
		if err = tx.Table(r.TableName()).
			Where("id=?", r.ParentId).
			Update("last_child_group_code", strings.TrimLeft(groupCodeV, parentGroupCode)).
			Error; err != nil {
			return
		}
		if err = tx.Table(r.TableName()).
			Where("id=?", r.Id).
			Update("group_code", groupCodeV).
			Error; err != nil {
			return
		}
	}
	return
}
func (r AdminGroup) AfterCreate(tx *gorm.DB) (err error) {
	if r.GroupCode == "" {
		var (
			groupCodeV      string
			parentGroupCode string
		)
		if groupCodeV, parentGroupCode, err = r.getGroupCode(tx); err != nil {
			return
		}
		if err = tx.Table(r.TableName()).
			Where("id=?", r.ParentId).
			Update("last_child_group_code", strings.TrimLeft(groupCodeV, parentGroupCode)).Error; err != nil {
			return
		}
		if err = tx.Model(r).Where("id=?", r.Id).
			Update("group_code", groupCodeV).Error; err != nil {
			return
		}
	}
	return
}

func (r *AdminGroup) initGroupCode(tx *gorm.DB) (err error) {
	var parentAll = make([]AdminGroup, 0, 15)
	r.OrgAllParent(tx, r.ParentId, &parentAll)
	var (
		parentGroup AdminGroup
	)
	for _, group := range parentAll {
		if group.GroupCode != "" {
			continue
		}

		if group.ParentId == 0 {
			group.GroupCode = (&groupCode{
				Prefix:     GroupCodePrefix,
				CodeLength: GroupCodeTopStepLength,
			}).GetCode()
			if err = tx.Table(group.TableName()).
				Where("id = ?", group.Id).
				Limit(1).
				Updates(map[string]interface{}{"group_code": group.GroupCode}).
				Error; err != nil {
				return
			}
			parentGroup = group
			continue
		} else if group.Id == 0 {
			var tList []AdminGroup
			if err = tx.Table(group.TableName()).
				Where("id = ?", group.Id).
				Limit(1).
				Find(&tList).
				Error; err != nil {
				return
			}
			if len(tList) == 0 {
				err = fmt.Errorf("数据异常（model:adming_group）")
				return
			}
			parentGroup = tList[0]
		}

		if group.GroupCode != "" {
			parentGroup = group
			continue
		}
		group.GroupCode = (&groupCode{
			Prefix:     parentGroup.GroupCode,
			FromString: parentGroup.LastChildGroupCode,
			CodeLength: GroupCodeStepLength,
		}).GetCode()
		if err = tx.Table(group.TableName()).
			Where("id = ?", group.Id).
			Limit(1).
			Updates(map[string]interface{}{"group_code": group.GroupCode}).
			Error; err != nil {
			return
		}
		if err = tx.Table(r.TableName()).
			Where("id=?", parentGroup.Id).
			Update("last_child_group_code", strings.TrimLeft(group.GroupCode, parentGroup.GroupCode)).Error; err != nil {
			return
		}

		parentGroup = group
	}

	return
}
func (r *AdminGroup) OrgAllParent(tx *gorm.DB, parentId int64, parentAll *[]AdminGroup) (err error) {
	if parentId == 0 {
		return
	}
	var m AdminGroup
	var m1 []AdminGroup
	if err = tx.Table(m.TableName()).
		Where("id = ?", parentId).
		Limit(1).
		Find(&m1).
		Error; err != nil {
		return
	}
	if len(m1) > 0 {
		*parentAll = append(m1, *parentAll...)
		err = r.OrgAllParent(tx, m1[0].ParentId, parentAll)
	}
	return
}

func (r *AdminGroup) getGroupCode(tx *gorm.DB) (res string, parentGroupCode string, err error) {
	if err = r.initGroupCode(tx); err != nil {
		return
	}
	var fromString string
	var dt []AdminGroup
	if err = tx.Model(r).
		Where("id = ?", r.ParentId).
		Limit(1).
		Find(&dt).Error; err != nil {
		return
	}
	if len(dt) > 0 {
		fromString = dt[0].LastChildGroupCode
		res = (&groupCode{
			Prefix:     dt[0].GroupCode,
			FromString: fromString,
		}).GetCode()
	}

	return

}

type groupCode struct {
	Prefix     string `json:"prefix"`
	FromString string `json:"from_string"`
	CodeLength int    `json:"code_length"`
}

func (r *groupCode) GetCode() (res string) {
	r.CodeLength = GroupCodeStepLength
	if r.Prefix == "" || r.Prefix == GroupCodePrefix {
		r.Prefix = GroupCodePrefix
		r.CodeLength = GroupCodeTopStepLength
	}

	res = r.Prefix + r.orgString(r.CodeLength)
	return
}

func (r *groupCode) orgString(codeLength int) (rs string) {
	s := []byte(`A0BCDE12FGHIJKLMNOP34QRST56UVW7XYZ89`)
	code := make([]byte, 0, codeLength)
	var s1 []byte
	fromCode := r.flipByte([]byte(r.FromString))
	var l = len(s)
	var nextNeedAdd bool
	var loc int
	for i := 0; i < codeLength; i++ {
		s1 = append(s[i+codeLength:], s[:i+codeLength]...)
		if r.FromString == "" {
			code = append(code, s1[0])
			continue
		}
		if i == 0 {
			if loc = r.getIndex(fromCode[i], s1) + 1; loc == l {
				nextNeedAdd = true
				loc = 0
			}
			code = append(code, s1[loc])
			continue
		}
		if nextNeedAdd {
			nextNeedAdd = false
			if loc = r.getIndex(fromCode[i], s1) + 1; loc == l {
				nextNeedAdd = true
				loc = 0
			}
			code = append(code, s1[loc])
			continue
		}
		code = append(code, fromCode[i])

	}
	rs = string(r.flipByte(code))
	return
}
func (r *groupCode) flipByte(arg []byte) (res []byte) {
	l := len(arg)
	res = make([]byte, l)
	for i := 0; i < l; i++ {
		res[l-i-1] = arg[i]
	}
	return
}
func (r *groupCode) getIndex(d byte, f []byte) (res int) {
	res = -1
	for k, value := range f {
		if d == value {
			res = k
			break
		}
	}
	return
}
