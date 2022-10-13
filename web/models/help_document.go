package models

import (
	"fmt"
	"github.com/juetun/base-wrapper/lib/base"
)

type HelpDocument struct {
	Id        int64            `gorm:"column:id;primary_key" json:"id" form:"id"`
	Label     string           `gorm:"column:label;not null;type:varchar(150);default:'';comment:帮助文档标题" json:"label"`
	Desc      string           `gorm:"column:desc;not null;type:varchar(150);default:'';comment:简单描述" json:"desc"`
	PKey      string           `gorm:"column:p_key;uniqueIndex:idx_pk,priority:1;not null;type:varchar(150);default:'';comment:唯一的key" json:"p_key"`
	Content   string           `gorm:"column:content;not null;type:text;comment:文档内容" json:"content"`
	CreatedAt base.TimeNormal  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"-" `
	UpdatedAt base.TimeNormal  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"-" `
	DeletedAt *base.TimeNormal `gorm:"column:deleted_at" json:"-"`
}

func (r *HelpDocument) GetTableComment() (res string) {
	res = "帮助文档"
	return
}

func (r *HelpDocument) DefaultBeforeAdd() (err error) {
	if r.PKey == "" {
		err = fmt.Errorf("请填写文档的唯一KEY")
		return
	}
	return
}
func (r *HelpDocument) TableName() string {
	return fmt.Sprintf("%shelp_document", TablePrefix)
}
