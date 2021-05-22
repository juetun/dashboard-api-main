/**
* @Author:changjiang
* @Description:
* @File:service
* @Version: 1.0.0
* @Date 2021/5/22 11:18 下午
 */
package wrappers

import (
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/common/response"
)

type (
	ArgServiceList struct {
		app_obj.JwtUserMessage
		response.BaseQuery
	}
	ResultServiceList struct {
		*response.Pager
	}
	ArgServiceEdit struct {
		Id        int    `json:"id" form:"id"`
		Name      string `json:"name" form:"name"`
		UniqueKey string `json:"unique_key" form:"column:unique_key"`
		Port      int    `json:"port"  form:"column:port"`
		Desc      string `json:"desc" form:"column:desc"`
		IsStop    uint8  `json:"is_stop" form:"column:is_stop"`
	}
	ResultServiceEdit struct {
		Result bool `json:"result"`
	}
)
