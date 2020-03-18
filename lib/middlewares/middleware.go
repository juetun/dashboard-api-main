package middlewares

import "github.com/gin-gonic/gin"

var MiddleWareComponent = []gin.HandlerFunc{}

func LoadMiddleWare() {

	cors := CORS(CORSOptions{Origin: "*",}) //跨域控制
	auth := Auth()                          //用户信息验证 token反解
	//webApp.Use(m.RequestID(m.RequestIDOptions{AllowSetting: true}))
	//webApp.Use(ginutil.Recovery(recoverHandler))
	//webApp.Use(m2.CheckExist())
	//webApp.Static("/static/uploads/images/", "./static/uploads/images/")
	MiddleWareComponent = append(MiddleWareComponent, cors, auth)
}
