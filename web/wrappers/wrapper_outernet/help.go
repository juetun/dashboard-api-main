package wrapper_outernet

import "github.com/gin-gonic/gin"

type (
	ArgTree struct {
		TopId     int64  `json:"top_id" form:"top_id"`
		BizCode   string `json:"biz_code" form:"biz_code"`
		CurrentId int64  `json:"current_id" form:"current_id"`
	}
	ResultTree struct {
		Data []*ResultFormPage `json:"data"`
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
)

func (r *ArgTree) Default(context *gin.Context) (err error) {

	return
}
