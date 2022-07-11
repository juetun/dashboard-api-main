package wrapper_outernet

import "github.com/gin-gonic/gin"

type (
	ArgTree struct {
		TopId     int64  `json:"top_id" form:"top_id"`
		BizCode   string `json:"biz_code" form:"biz_code"`
		CurrentId int64  `json:"current_id" form:"current_id"`
		DocKey    string `json:"doc_key" form:"doc_key"`
	}
	ResultTree struct {
		Data []*ResultFormPage `json:"data"`
		DocContent string `json:"doc_content"`
	}
	ResultFormPage struct {
		Id         int64             `json:"id"`
		Title      string            `json:"title"`
		Expand     bool              `json:"expand,omitempty"`
		DocKey     string            `json:"doc_key,omitempty"`
 		IsLeafNode uint8             `json:"is_leaf_node,omitempty"`
		Children   []*ResultFormPage `json:"children,omitempty"`
	}
)

func (r *ArgTree) Default(context *gin.Context) (err error) {

	return
}
