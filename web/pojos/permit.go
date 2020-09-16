/**
* @Author:changjiang
* @Description:
* @File:permit
* @Version: 1.0.0
* @Date 2020/9/16 12:22 上午
 */
package pojos

import (
	"github.com/juetun/base-wrapper/lib/app_obj"
	"github.com/juetun/dashboard-api-main/web/models"
)

type ArgPermitMenu struct {
	app_obj.JwtUserMessage
}
type ResultPermitMenu struct {
	Menu []models.AdminMenu
}

type ArgFlag struct {
	app_obj.JwtUserMessage
}

type ResultFlag struct {
}
