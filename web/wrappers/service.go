// Package wrappers
/**
* @Author:changjiang
* @Description:
* @File:service
* @Version: 1.0.0
* @Date 2021/5/22 11:18 下午
 */
package wrappers

import (
	"fmt"
	"github.com/juetun/base-wrapper/lib/base"
	"strconv"
	"time"

	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/library/common/app_param"
)

type (
	ImportMenu struct {
		AppName string `json:"app_name"`
	}
	ArgDetail struct {
		app_param.RequestUser
		Id int `json:"id"`
	}
	ResultDetail struct {
		models.AdminApp
	}
	ArgServiceList struct {
		app_param.RequestUser
		response.PageQuery
		Name       string   `json:"name" form:"name"`
		Id         int      `json:"id" form:"id" `
		UniqueKey  string   `json:"unique_key" form:"unique_key"`
		UniqueKeys []string `json:"-" form:"-"`
		Port       int      `json:"port"  form:"port"`
		Desc       string   `json:"desc" form:"desc"`
		IsStop     uint8    `json:"is_stop" form:"is_stop"`
	}
	ResultServiceList struct {
		response.Pager
	}
	AdminApp struct {
		Id            int       `json:"id" gorm:"column:id;primary_key" `
		UniqueKey     string    `json:"unique_key" gorm:"column:unique_key"`
		Port          int       `json:"port"  gorm:"column:port"`
		Name          string    `json:"name" gorm:"column:name"`
		Desc          string    `json:"desc" gorm:"column:desc"`
		IsStop        uint8     `json:"is_stop" gorm:"column:is_stop"`
		IsStopName    string    `json:"is_stop_name" gorm:"-"`
		Expanded      bool      `json:"_expanded" gorm:"-"`
		DisableExpand bool      `json:"_disableExpand" gorm:"-"` //禁用当前行的展开功能
		CreatedAt     time.Time `json:"created_at" gorm:"column:created_at" `
		UpdatedAt     time.Time `json:"updated_at" gorm:"column:updated_at" `
	}
	// ArgServiceEdit {"id":1,"unique_key":"app-user","port":80,"name":"用户","desc":"","is_stop":1}
	ArgServiceEdit struct {
		Name       string            `json:"name" form:"name"`
		UniqueKey  string            `json:"unique_key" form:"unique_key"`
		HostConfig map[string]string `json:"hosts" form:"hosts"`
		Port       int               `json:"port"  form:"port"`
		Desc       string            `json:"desc" form:"desc"`
		IsStop     int               `json:"is_stop" form:"is_stop"`
		Id         int               `json:"id" form:"id"`
	}

	ResultServiceEdit struct {
		Result bool `json:"result"`
	}
)

func (r *ArgServiceList) Default(ctx *base.Context) (err error) {

	return
}

func (r *ArgDetail) Default(ctx *base.Context) (err error) {
	if r.Id, err = strconv.Atoi(ctx.GinContext.Params.ByName("id")); err != nil {
		err = fmt.Errorf("参数格式不正确")
		return
	}
	return
}

func (r *ArgServiceEdit) Default(c *base.Context) (err error) {
	if r.IsStop < 0 || r.IsStop > 100 {
		err = fmt.Errorf("is_stop格式不正确")
		return
	}
	return
}
