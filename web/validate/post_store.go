/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2019-01-12
 * Time: 19:20
 */
package validate

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/dashboard-api-main/web/pojos"
)

type PostStoreV struct {
}

func (pv *PostStoreV) MyValidate() gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := common.NewGin(c)
		var json pojos.PostStore
		// 接收各种参数
		if err := c.ShouldBindJSON(&json); err != nil {
			appG.Response(400001000, err.Error())
			c.Abort()
			return
		}

		reqValidate := &PostStore{
			Title:   json.Title,
			Tags:    json.Tags,
			Summary: json.Summary,
		}
		if b := appG.Validate(reqValidate); !b {
			c.Abort()
			return
		}
		c.Set("json", json)
		c.Next()
	}
}

type PostStore struct {
	Title string `valid:"Required"`
	Tags  []int
	// Category int `valid:Required`
	Summary string `valid:"Required;"`
}

func (c *PostStore) Message() map[string]common.ValidationMessage {
	return map[string]common.ValidationMessage{
		"Title.Required.":   {Code: 401000000, Message: "请输入标题"},
		"Summary.Required.": {Code: 401000003, Message: "请输入摘要"},
	}
}
