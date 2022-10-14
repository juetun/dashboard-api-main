package wrapper_admin

import (
	"fmt"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/base-wrapper/lib/utils"
	"github.com/juetun/dashboard-api-main/web/models"
	"time"
)

type (
	ArgHelpTree struct {
		TopId     int64  `json:"top_id" form:"top_id"`
		BizCode   string `json:"biz_code" form:"biz_code"`
		CurrentId int64  `json:"current_id" form:"current_id"`
	}
	ResultHelpTree struct {
		List []*ResultFormPage         `json:"list"`
		Menu []*ResultHelpTreeItemMenu `json:"menu"` // 一级系统权限列表
	}

	ResultHelpTreeItem struct {
		Child []*ResultHelpTreeItem `json:"child,omitempty"`
		ResultHelpTreeItemMenu
	}

	ResultHelpTreeItemMenu struct {
		models.HelpDocumentRelate
		Active bool `json:"active"`
	}

	ResultFormPage struct {
		Id         int64             `json:"id"`
		Label      string            `json:"label"`
		Expand     bool              `json:"expand"`
		DocKey     string            `json:"doc_key"`
		Display    uint8             `json:"display"`
		IsLeafNode uint8             `json:"is_leaf_node"`
		BizCode    string            `json:"biz_code"`
		ParentId   int64             `json:"parent_id"`
		Children   []*ResultFormPage `json:"children,omitempty"`
	}
	ArgTreeEditNode struct {
		Id            int64           `json:"id" form:"id"`
		BizCode       string          `json:"biz_code" form:"biz_code"`
		Label         string          `json:"label" form:"label"`
		Display       uint8           `json:"display" form:"display"`
		ParentId      int64           `json:"parent_id" form:"parent_id"`
		IsLeafNode    uint8           `json:"is_leaf_node" form:"is_leaf_node"`
		DocKey        string          `json:"doc_key" form:"doc_key"`
		IsBizCodeEdit uint8           `json:"is_biz_code_edit" form:"is_biz_code_edit"`
		TimeNow       base.TimeNormal `json:"-" form:"-"`
	}
	ResultTreeEditNode struct {
		Result bool `json:"result"`
	}
	ArgHelpList struct {
		response.PageQuery
		PKey    string    `json:"p_key" form:"p_key"`
		TimeNow time.Time `json:"-" form:"-"`
	}
	ResultHelpList struct {
		response.Pager
	}
	ResultHelpListItem struct {
		Id        int64  `json:"id"`
		Label     string `json:"label"`
		Desc      string `json:"desc"`
		PKey      string `json:"p_key"`
		Content   string `json:"content"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}
	ArgHelpDetail struct {
		Id   int64  `json:"id" form:"id"`
		PKey string `json:"p_key" form:"p_key"`
	}

	ResultHelpDetail struct {
		Id      int64  `json:"id"`
		Label   string `json:"label"`
		Desc    string `json:"desc"`
		PKey    string `json:"p_key"`
		Content string `json:"content"`
	}

	ArgHelpEdit struct {
		Id      int64           `json:"id" form:"id"`
		Label   string          `json:"label" form:"label"`
		Desc    string          `json:"desc" form:"desc"`
		PKey    string          `json:"p_key" form:"p_key"`
		Content string          `json:"content" form:"content"`
		TimeNow base.TimeNormal `json:"-" form:"-"`
	}

	ResultHelpEdit struct {
		Result bool `json:"result"`
	}
)

func (r *ResultHelpTreeItemMenu) SetHelpDocumentRelate(relate *models.HelpDocumentRelate) {
	r.HelpDocumentRelate = *relate
}

func NewResultFormPage() (res *ResultFormPage) {
	return &ResultFormPage{}

}
func (r *ResultFormPage) SetResultHelpTreeItem(item *ResultHelpTreeItem, currentId int64) (res *ResultFormPage) {
	res = r
	r.Label = item.Label
	r.Id = item.Id
	r.DocKey = item.DocKey
	r.Display = item.Display
	r.IsLeafNode = item.IsLeafNode
	r.BizCode = item.BizCode
	r.ParentId = item.ParentId
	if item.Id == currentId {
		r.Expand = true
	}
	return
}

func NewResultHelpTree() (res *ResultHelpTree) {
	res = &ResultHelpTree{
		List: []*ResultFormPage{},
		Menu: []*ResultHelpTreeItemMenu{},
	}
	return
}
func (r *ResultHelpListItem) SetHelpDocument(document *models.HelpDocument, timeNow time.Time) {
	r.Id = document.Id
	r.Label = document.Label
	r.Desc = document.Desc
	r.PKey = document.PKey
	r.Content = document.Content
	r.CreatedAt, _, _ = utils.DateTimeDiff(timeNow, document.CreatedAt.Time, utils.DateTimeDashboard)
	r.UpdatedAt, _, _ = utils.DateTimeDiff(timeNow, document.UpdatedAt.Time, utils.DateTimeDashboard)
	return
}

func (r *ArgHelpTree) Default(context *base.Context) (err error) {

	return
}

func (r *ArgTreeEditNode) Default(context *base.Context) (err error) {
	r.TimeNow = base.GetNowTimeNormal()
	if r.Label == "" {
		err = fmt.Errorf("请填写标题")
		return
	}
	if r.IsBizCodeEdit == 0 {
		if r.BizCode == "" {
			err = fmt.Errorf("请填写业务编码")
			return
		}
		if r.IsLeafNode == models.HelpDocumentRelateIsLeafNodeYes { //如果是叶子节点
			if r.DocKey == "" {
				err = fmt.Errorf("叶子节点必须填写KEY")
				return
			}
		}
	} else {
		r.BizCode = r.DocKey
		r.IsLeafNode = models.HelpDocumentRelateIsLeafNodeNo
	}

	return
}

func (r *ResultHelpDetail) ParseFromHelpDoc(document *models.HelpDocument) {
	r.Id = document.Id
	r.PKey = document.PKey
	r.Label = document.Label
	r.Desc = document.Desc
	r.Content = document.Content
	return
}

func (r *ArgHelpList) Default(c *base.Context) (err error) {
	r.TimeNow = time.Now()
	return
}

func (r *ArgHelpDetail) Default(c *base.Context) (err error) {
	if r.Id == 0 && r.PKey == "" {
		err = fmt.Errorf("请选择你要查看的信息")
		return
	}
	return
}

func (r *ArgHelpEdit) Default(context *base.Context) (err error) {
	if r.TimeNow.IsZero() {
		r.TimeNow = base.GetNowTimeNormal()
	}
	var maxLength = 150
	if utils.StringUtf8Length(r.PKey) > maxLength {
		err = fmt.Errorf("p_key不超过%d个字符", maxLength)
		return
	}
	return
}
