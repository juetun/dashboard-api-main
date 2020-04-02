/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2019-05-07
 * Time: 23:30
 */
package validate

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juetun/app-dashboard/lib/base"
	"github.com/juetun/app-dashboard/lib/common"
	"github.com/juetun/app-dashboard/web/pojos"
)

type SystemUpdateV struct {
}

func (sv *SystemUpdateV) MyValidate() gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := common.NewGin(c)
		var json pojos.ConsoleSystem
		// 接收各种参数
		if err := c.ShouldBindJSON(&json); err != nil {
			emsg := "参数验证失败" + err.Error()
			c.JSON(http.StatusOK, base.Result{Code: 400001000, Data: nil, Msg: emsg})
			c.Abort()
			return
		}

		reqValidate := &SystemUpdate{
			Title:        json.Title,
			Keywords:     json.Keywords,
			Description:  json.Description,
			RecordNumber: json.RecordNumber,
			Theme:        json.Theme,
		}
		if b := appG.Validate(reqValidate); !b {
			c.Abort()
			return
		}
		c.Set("json", json)
		c.Next()
	}
}

type SystemUpdate struct {
	Title        string `valid:"Required;MaxSize(100)"`
	Keywords     string `valid:"Required;MaxSize(100)"`
	Description  string `valid:"Required;MaxSize(250)"`
	RecordNumber string `valid:"Required;MaxSize(50)"`
	Theme        int    `valid:"Required;"`
}

func (c *SystemUpdate) Message() map[string]common.ValidationMessage {
	return map[string]common.ValidationMessage{
		"Title.Required.":        {Code: 405000001, Message: "请输入网站title"},
		"Title.MaxSize.":         {Code: 405000002, Message: "网站title不超过100个字符"},
		"Keywords.Required.":     {Code: 405000003, Message: "请输入网站关键字"},
		"Keywords.MaxSize.":      {Code: 405000004, Message: "网站关键字不超过100个字符"},
		"Description.Required.":  {Code: 405000005, Message: "请输入网站描述"},
		"Description.MaxSize.":   {Code: 405000006, Message: "网站描述不超过250个字符"},
		"RecordNumber.Required.": {Code: 405000007, Message: "请输入备案号"},
		"RecordNumber.MaxSize.":  {Code: 405000008, Message: "请输入备案号不超过50个字符"},
		"Theme.Required.":        {Code: 405000009, Message: "请输入主题"},
	}
}
