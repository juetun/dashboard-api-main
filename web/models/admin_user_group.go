/**
* @Author:changjiang
* @Description:
* @File:admin_user_group
* @Version: 1.0.0
* @Date 2020/9/16 10:38 下午
 */
package models

type AdminUserGroup struct {
	Id      int    `gorm:"primary_key" json:"id"`
	GroupId int    `json:"group_id" gorm:"group_id"`
	UserHid string `json:"user_hid"  gorm:"group_id"`
	IsDel   int    `json:"-" gorm:"is_del"`
}

func (r *AdminUserGroup) TableName() string {
	return "admin_user_group"
}
