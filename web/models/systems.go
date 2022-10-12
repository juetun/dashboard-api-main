package models

import (
	"fmt"
	"github.com/juetun/base-wrapper/lib/base"
)

type ZBaseSys struct {
	base.Model
	Theme        int    `gorm:"column:theme;default:0;type:int(10) COLLATE utf8mb4_bin;comment:主题" json:"theme"`
	Title        string `gorm:"column:title;default:'';type:varchar(100) COLLATE utf8mb4_general_ci;comment:标题" json:"title"`
	Keywords     string `gorm:"column:keywords;default:'';type:varchar(300) COLLATE utf8mb4_general_ci;comment:网站关键词" json:"keywords"`
	Description  string `gorm:"column:description;default:'';type:varchar(5000) COLLATE utf8mb4_general_ci;comment:网站描述" json:"description"`
	RecordNumber string `gorm:"column:record_number;default:'';type:varchar(100) COLLATE utf8mb4_bin;comment:网站备案号" json:"record_number"`
}

func (r *ZBaseSys) GetTableComment() (res string) {
	return "系统设置"
}

func (r *ZBaseSys) TableName() string {
	return fmt.Sprintf("%sbase_sys", TablePrefix)
}
