package models

import (
	"fmt"
	"github.com/juetun/base-wrapper/lib/base"
)

var (
	SliceOperateLogModule = base.ModelItemOptions{}
)

type OperateLog struct {
	Id          int64           `gorm:"column:id;primary_key" json:"id" form:"id"`
	UserHid     int64           `json:"user_hid" gorm:"column:user_hid;index:idx_user_hid,priority:1;not null;type:bigint(20);default:0;comment:用户id"`
	Module      string          `json:"module" gorm:"column:module;default:'';not null;type:varchar(100);comment:所属模块" form:"module"`
	DataType    string          `json:"data_type"  gorm:"column:data_type;default:'';not null;type:varchar(100);comment:数据类型" form:"data_type"` //数据类型
	DataId      string          `json:"data_id"  gorm:"column:data_id;default:'';not null;type:varchar(100);comment:数据ID" form:"data_id"`       //数据ID
	Description string          `gorm:"column:description;type:text;not null;comment:日志内容;" json:"description"`
	CreatedAt   base.TimeNormal `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at" `
}

func (r *OperateLog) TableName() (res string) {
	return fmt.Sprintf("%soperate_log", TablePrefix)
}

func (r *OperateLog) GetTableComment() (res string) {

	return "客服后台操作日志"
}

func (r *OperateLog) ParseModule() (res string, err error) {
	var mapModule map[string]string
	if mapModule, err = SliceOperateLogModule.GetMapAsKeyString(); err != nil {
		return
	}
	if tmp, ok := mapModule[r.Module]; ok {
		res = tmp
		return
	}
	res = r.Module
	return
}
