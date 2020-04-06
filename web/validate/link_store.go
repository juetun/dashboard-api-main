/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2019-05-08
 * Time: 23:00
 */
package validate

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/app-dashboard/web/pojos"
)

type LinkStoreV struct {
}

func (lv *LinkStoreV) MyValidate() gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := common.NewGin(c)
		var json pojos.LinkStore
		if err := c.ShouldBindJSON(&json); err != nil {
			appG.Response(400001000, err.Error())
			c.Abort()
			return
		}

		reqValidate := &LinkStore{
			Name:  json.Name,
			Link:  json.Link,
			Order: json.Order,
		}
		if b := appG.Validate(reqValidate); !b {
			c.Abort()
			return
		}
		c.Set("json", json)
		c.Next()
	}
}

type LinkStore struct {
	Name  string `valid:"Required;MaxSize(100)"`
	Link  string `valid:"Required;MaxSize(1000)"`
	Order int    `valid:"Min(0);Max(1000000000)"`
}

func (c *LinkStore) Message() map[string]common.ValidationMessage {
	return map[string]common.ValidationMessage{
		"Name.Required.": {Code: 406000000, Message: "请输入外链名称"},
		"Name.MaxSize.":  {Code: 406000001, Message: "请输入外链名称不超过100个字符"},
		"Link.Required.": {Code: 406000002, Message: "请输入连接地址"},
		"Link.MaxSize.":  {Code: 406000003, Message: "链接地址不超过1000个字符"},
		"Order.Min.":     {Code: 406000004, Message: "排序值不能小于0"},
		"Order.Max.":     {Code: 406000004, Message: "排序值不能大于1000000000"},
	}
}
