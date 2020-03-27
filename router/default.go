package router

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/app-dashboard/web/controllers/web_impl"
	"html/template"
	"time"
)

func init() {
	HandleFunc = append(HandleFunc, func(r *gin.Engine, urlPrefix string) {
		web := web_impl.NewControllerWeb()

		h := r.Group(urlPrefix)
		{
			r.SetFuncMap(template.FuncMap{
				"rem":    rem,
				"MDate":  mDate,
				"MDate2": mDate2,
			})
			r.LoadHTMLGlob("template/home/*.tmpl")

			r.Static(urlPrefix +"/static/home", "./static/home")
			h.GET(urlPrefix +"/", web.Index)
			h.GET(urlPrefix +"/categories/:name", web.IndexCate)
			h.GET(urlPrefix +"/tags/:name", web.IndexTag)
			h.GET(urlPrefix +"/detail/:id", web.Detail)
			h.GET(urlPrefix +"/archives", web.Archives)
			h.GET(urlPrefix +"/404", web.NoFound)
		}
	}, )
}

func rem(divisor int) bool {
	if (divisor+1)%4 == 0 {
		return true
	} else {
		return false
	}
}

func mDate(times time.Time) string {
	return times.Format("2006-01-02 15:04:05")
}

func mDate2(times time.Time) string {
	return times.Format("01-02")
}
