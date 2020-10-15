/**
* @Author:changjiang
* @Description:
* @File:admin_user_group_permit
* @Version: 1.0.0
* @Date 2020/9/16 10:38 下午
 */
package models

type AdminUserGroupPermit struct {
	Id       int    `gorm:"primary_key" json:"id"`
	GroupId  int    `json:"group_id" gorm:"group_id" `
	PathType string `json:"path_type" gorm:"path_type"`
	MenuId   int    `json:"menu_id" gorm:"menu_id"`
	IsDel    int    `json:"-" gorm:"is_del"`
}

func (r *AdminUserGroupPermit) TableName() string {
	return "admin_user_group_permit"
}
