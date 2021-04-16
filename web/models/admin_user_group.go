/**
* @Author:changjiang
* @Description:
* @File:admin_user_group
* @Version: 1.0.0
* @Date 2020/9/16 10:38 下午
 */
package models

type AdminUserGroup struct {
	Id      int    `gorm:"column:primary_key" json:"id"`
	GroupId int    `json:"group_id" gorm:"column:group_id"`
	UserHid string `json:"user_hid"  gorm:"column:user_hid"`
	IsDel   int    `json:"-" gorm:"column:is_del"`
}

func (r *AdminUserGroup) TableName() string {
	return "admin_user_group"
}
