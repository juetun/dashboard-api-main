/**
* @Author:changjiang
* @Description:
* @File:admin_group
* @Version: 1.0.0
* @Date 2020/9/20 6:40 下午
 */
package models

type AdminGroup struct {
	Id       int    `gorm:"primary_key" json:"id"`
	Name     string `json:"name" gorm:"column:name"`
 	IsDel    int    `json:"-" gorm:"column:is_del"`
}

func (r *AdminGroup) TableName() string {
	return "admin_group"
}
