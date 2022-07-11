package wrapper_admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/base-wrapper/lib/utils"
	"github.com/juetun/dashboard-api-main/web/models"
)

type (
	ArgHelpTree struct {
		TopId     int64  `json:"top_id" form:"top_id"`
		BizCode   string `json:"biz_code" form:"biz_code"`
		CurrentId int64  `json:"current_id" form:"current_id"`
	}
	ResultHelpTree     []*ResultFormPage
	ResultHelpTreeItem struct {
		models.HelpDocumentRelate
 		Child []*ResultHelpTreeItem `json:"child,omitempty"`
	}
	ResultFormPage struct {
		Id         int64             `json:"id"`
		Title      string            `json:"title"`
		Expand     bool              `json:"expand"`
		DocKey     string            `json:"doc_key"`
		Display    uint8             `json:"display"`
		IsLeafNode uint8             `json:"is_leaf_node"`
		Children   []*ResultFormPage `json:"children,omitempty"`
	}
	ArgTreeEditNode struct {
		Id         int64           `json:"id" form:"id"`
		BizCode    string          `json:"biz_code" form:"biz_code"`
		Label      string          `json:"label" form:"label"`
		Display    uint8           `json:"display" form:"display"`
		ParentId   int64           `json:"parent_id" form:"parent_id"`
		IsLeafNode uint8           `json:"is_leaf_node" form:"is_leaf_node"`
		DocKey     string          `json:"doc_key" form:"doc_key"`
		TimeNow    base.TimeNormal `json:"-" form:"-"`
	}
	ResultTreeEditNode struct {
		Result bool `json:"result"`
	}
	ArgHelpList struct {
		response.PageQuery
		PKey string `json:"p_key" form:"p_key"`
	}
	ResultHelpList struct {
		*response.Pager
	}
	ArgHelpDetail struct {
		Id int64 `json:"id" form:"id"`
	}

	ResultHelpDetail struct {
		Id      int64  `json:"id"`
		PKey    string `json:"p_key"`
		Content string `json:"content"`
	}

	ArgHelpEdit struct {
		Id      int64           `json:"id" form:"id"`
		PKey    string          `json:"p_key" form:"p_key"`
		Content string          `json:"content" form:"content"`
		TimeNow base.TimeNormal `json:"-" form:"-"`
	}

	ResultHelpEdit struct {
		Result bool `json:"result"`
	}
)

func (r *ArgHelpTree) Default(context *gin.Context) (err error) {

	return
}

func (r *ArgTreeEditNode) Default(context *gin.Context) (err error) {
	r.TimeNow = base.GetNowTimeNormal()
	if r.BizCode == "" {
		err = fmt.Errorf("请填写业务编码")
		return
	}
	if r.Label == "" {
		err = fmt.Errorf("请填写标题")
		return
	}
	return
}

func (r *ResultHelpDetail) ParseFromHelpDoc(document *models.HelpDocument) {
	r.Id = document.Id
	r.PKey = document.PKey
	r.Content = document.Content
	return
}

func (r *ArgHelpList) Default(c *gin.Context) (err error) {

	return
}

func (r *ArgHelpDetail) Default(c *gin.Context) (err error) {
	if r.Id == 0 {
		err = fmt.Errorf("请选择你要查看的信息")
		return
	}
	return
}

func (r *ArgHelpEdit) Default(context *gin.Context) (err error) {
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
