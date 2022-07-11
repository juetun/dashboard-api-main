package models

import (
	"fmt"
	"github.com/juetun/base-wrapper/lib/base"
)

const (
	HelpDocumentRelateIsLeafNodeNo  uint8 = iota //不是叶子节点
	HelpDocumentRelateIsLeafNodeYes              //是叶子节点
)

const (
	DisplayNo uint8 = iota
	DisplayYes
)

var (
	SliceHelpDocumentRelateDisplay = base.ModelItemOptions{
		{
			Value: DisplayYes,
			Label: "展示",
		},
		{
			Value: DisplayNo,
			Label: "不展示",
		},
	}
	SliceHelpDocumentRelateIsLeafNode = base.ModelItemOptions{
		{
			Value: HelpDocumentRelateIsLeafNodeNo,
			Label: "非叶子节点",
		},
		{
			Value: HelpDocumentRelateIsLeafNodeYes,
			Label: "叶子节点",
		},
	}
)

type HelpDocumentRelate struct {
	Id         int64            `gorm:"column:id;primary_key" json:"id" form:"id"`
	BizCode    string           `gorm:"column:biz_code;not null;uniqueIndex:idx_pk,priority:1;type:varchar(150);default:'';comment:业务场景码"json:"biz_code"`
	Display    uint8            `gorm:"column:display;default:1;not null;type:tinyint(2);comment:是否在列表页展示 1-展示 0-不展示" json:"display"`
	ParentId   int64            `gorm:"column:parent_id;not null;uniqueIndex:idx_pk,priority:2;type:bigint(20);default:0;comment:上级文档ID" json:"parent_id"`
	Label      string           `gorm:"column:label;not null;uniqueIndex:idx_pk,priority:3;type:varchar(150);default:'';comment:名称"json:"label"`
	IsLeafNode uint8            `gorm:"column:is_leaf_node;default:0;not null;type:tinyint(2);comment:是否叶子节点;1-是 0-不是" json:"is_leaf_node"`
	DocKey     string           `gorm:"column:p_key;not null;type:varchar(150);default:'';comment:文档唯一的key"json:"doc_key"`
	CreatedAt  base.TimeNormal  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"-" `
	UpdatedAt  base.TimeNormal  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"-" `
	DeletedAt  *base.TimeNormal `gorm:"column:deleted_at" json:"-"`
}

func (r *HelpDocumentRelate) GetTableComment() (res string) {
	res = "帮助文档关系描述"
	return
}

func (r *HelpDocumentRelate) DefaultBeforeAdd() (err error) {

	switch r.IsLeafNode {
	case HelpDocumentRelateIsLeafNodeYes: //如果是叶子节点
		r.DocKey = ""
	default:
		err = fmt.Errorf("没有选择叶子节点对应的文档")
	}
	return
}

func (r *HelpDocumentRelate) ParseIsLeafNode() (res string) {
	mapIsLeafNode, _ := SliceHelpDocumentRelateIsLeafNode.GetMapAsKeyUint8()
	var ok bool
	if res, ok = mapIsLeafNode[r.IsLeafNode]; !ok {
		res = fmt.Sprintf("未知类型(%d)", r.IsLeafNode)
	}
	return
}

func (r *HelpDocumentRelate) ParseDisplay() (res string) {
	mapDisplay, _ := SliceHelpDocumentRelateDisplay.GetMapAsKeyUint8()
	var ok bool
	if res, ok = mapDisplay[r.Display]; !ok {
		res = fmt.Sprintf("未知类型(%d)", r.Display)
	}
	return
}

func (r *HelpDocumentRelate) TableName() string {
	return fmt.Sprintf("%shelp_document_relate", TablePrefix)
}
