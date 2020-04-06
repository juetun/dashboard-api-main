/**
 * Created by GoLand.
 * User: zhu
 * Email: ylsc633@gmail.com
 * Date: 2019-04-23
 * Time: 17:46
 */
package validate

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/dashboard-api-main/web/pojos"
)

type TagStoreV struct {
}

func (tv *TagStoreV) MyValidate() gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := common.NewGin(c)
		var json pojos.TagStore
		// 接收各种参数
		if err := c.ShouldBindJSON(&json); err != nil {
			appG.Response(400001000, err.Error())
			c.Abort()
			return
		}

		reqValidate := &TagStore{
			Name:        json.Name,
			DisplayName: json.DisplayName,
			SeoDesc:     json.SeoDesc,
		}
		if b := appG.Validate(reqValidate); !b {
			c.Abort()
			return
		}
		c.Set("json", json)
		c.Next()
	}
}

type TagStore struct {
	Name        string `valid:"Required;MaxSize(100)"`
	DisplayName string `valid:"Required;MaxSize(100)"`
	SeoDesc     string `valid:"Required;MaxSize(250)"`
}

func (c *TagStore) Message() map[string]common.ValidationMessage {
	return map[string]common.ValidationMessage{
		"Name.Required.":        {Code: 403000000, Message: "请输入标签"},
		"Name.MaxSize.":         {Code: 403000001, Message: "标签不超过100个字符"},
		"DisplayName.Required.": {Code: 403000002, Message: "请输入标签别名"},
		"DisplayName.MaxSize.":  {Code: 403000003, Message: "标签别名不超过100个字符"},
		"SeoDesc.Required.":     {Code: 403000004, Message: "请输入seo描述"},
		"SeoDesc.MaxSize.":      {Code: 403000005, Message: "seo描述不超过250个字符"},
	}
}
