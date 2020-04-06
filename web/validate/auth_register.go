/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2019-05-11
 * Time: 00:43
 */
package validate

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/app-dashboard/web/pojos"
)

type AuthRegisterV struct {
}

func (r *AuthRegisterV) MyValidate() gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := common.NewGin(c)
		var json pojos.AuthRegister
		if err := c.ShouldBindJSON(&json); err != nil {
			// o.Response(http.StatusOK, 400001000, nil)
			appG.Response(400001000, err.Error())
			c.Abort()
			return
		}

		reqValidate := &AuthRegister{
			Email:      json.Email,
			Password:   json.Password,
			Captcha:    json.Captcha,
			CaptchaKey: json.CaptchaKey,
		}
		if b := appG.Validate(reqValidate); !b {
			c.Abort()
			return
		}
		c.Set("json", json)
		c.Next()
	}
}

type AuthRegister struct {
	Email      string `valid:"Required;Email"`
	Password   string `valid:"Required;MaxSize(30)"`
	Captcha    string `valid:"Required;MaxSize(5)"`
	CaptchaKey string `valid:"Required;MaxSize(30)"`
}

func (r *AuthRegister) Message() map[string]common.ValidationMessage {
	return map[string]common.ValidationMessage{
		"Email.Required.":      {Code: 407000000, Message: "请输入邮箱"},
		"Email.Email.":         {Code: 407000001, Message: "您输入的邮箱格式不正确"},
		"Password.Required.":   {Code: 407000002, Message: "请输入密码"},
		"Password.MaxSize.":    {Code: 407000003, Message: "密码不超过30个字符"},
		"Captcha.Required.":    {Code: 407000004, Message: "请输入验证码"},
		"Captcha.MaxSize.":     {Code: 407000005, Message: "验证码格式不正确"},
		"CaptchaKey.Required.": {Code: 407000006, Message: "数据异常，验证校验码"},
		"CaptchaKey.MaxSize.":  {Code: 407000007, Message: "校验参数不超过30个字符"},
	}

}
