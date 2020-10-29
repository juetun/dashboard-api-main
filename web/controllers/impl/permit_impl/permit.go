/**
* @Author:changjiang
* @Description:
* @File:permit
* @Version: 1.0.0
* @Date 2020/9/16 12:06 上午
 */
package permit_impl

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/controllers/inter"
	"github.com/juetun/dashboard-api-main/web/pojos"
	"github.com/juetun/dashboard-api-main/web/services"
)

type ControllerPermit struct {
	base.ControllerBase
}

func NewControllerPermit() inter.Permit {
	controller := &ControllerPermit{}
	controller.ControllerBase.Init()
	return controller
}

func (r *ControllerPermit) AdminUserAdd(c *gin.Context) {
	var arg pojos.ArgAdminUserAdd
	var err error
	err = c.Bind(&arg)
	if err != nil {
		r.Response(c, 500000001, nil, err.Error())
		return
	}
	arg.JwtUserMessage = r.GetUser(c)
	arg.Default()

	// 记录日志
	res, err := services.
		NewPermitService(base.GetControllerBaseContext(&r.ControllerBase, c)).
		AdminUserAdd(&arg)

	if err != nil {
		r.Response(c, 500000002, nil, err.Error())
		return
	}
	r.Response(c, 0, res)
}

func (r *ControllerPermit) AdminUserDelete(c *gin.Context) {
	var arg pojos.ArgAdminUserDelete
	var err error
	err = c.Bind(&arg)
	if err != nil {
		r.Response(c, 500000001, nil, err.Error())
		return
	}
	arg.JwtUserMessage = r.GetUser(c)
	arg.Default()

	if len(arg.IdString) == 0 {
		r.Response(c, 500000001, nil, "您没有选择要删除的用户")
		return
	}
	// 记录日志
	res, err := services.
		NewPermitService(base.GetControllerBaseContext(&r.ControllerBase, c)).
		AdminUserDelete(&arg)
	if err != nil {
		r.Response(c, 500000002, nil, err.Error())
		return
	}
	r.Response(c, 0, res)
}

func (r *ControllerPermit) AdminUserGroupRelease(c *gin.Context) {
	var arg pojos.ArgAdminUserGroupRelease
	var err error
	err = c.Bind(&arg)
	if err != nil {
		r.Response(c, 500000001, nil, err.Error())
		return
	}
	arg.JwtUserMessage = r.GetUser(c)
	arg.Default()

	// 记录日志
	srv := services.NewPermitService(base.GetControllerBaseContext(&r.ControllerBase, c))
	res, err := srv.AdminUserGroupRelease(&arg)

	if err != nil {
		r.Response(c, 500000002, nil, err.Error())
		return
	}
	r.Response(c, 0, res)
}

func (r *ControllerPermit) AdminUserGroupAdd(c *gin.Context) {
	var arg pojos.ArgAdminUserGroupAdd
	var err error
	err = c.Bind(&arg)
	if err != nil {
		r.Response(c, 500000001, nil, err.Error())
		return
	}
	arg.JwtUserMessage = r.GetUser(c)
	arg.Default()

	// 记录日志
	context := base.GetControllerBaseContext(&r.ControllerBase, c)
	context.Log.Logger.Infof("user:%+v", arg.JwtUserMessage)

	srv := services.NewPermitService(context)
	res, err := srv.AdminUserGroupAdd(&arg)

	if err != nil {
		r.Response(c, 500000002, nil, err.Error())
		return
	}
	r.Response(c, 0, res)
}

func (r *ControllerPermit) AdminGroupDelete(c *gin.Context) {
	var arg pojos.ArgAdminGroupDelete
	var err error
	err = c.Bind(&arg)
	if err != nil {
		r.Response(c, 500000001, nil, err.Error())
		return
	}
	arg.JwtUserMessage = r.GetUser(c)
	arg.Default()

	// 记录日志
	context := base.GetControllerBaseContext(&r.ControllerBase, c)
	context.Log.Logger.Infof("user:%+v", arg.JwtUserMessage)

	srv := services.NewPermitService(context)
	res, err := srv.AdminGroupDelete(&arg)

	if err != nil {
		r.Response(c, 500000002, nil, err.Error())
		return
	}
	r.Response(c, 0, res)
}
func (r *ControllerPermit) AdminGroupEdit(c *gin.Context) {
	var arg pojos.ArgAdminGroupEdit
	var err error
	err = c.Bind(&arg)
	if err != nil {
		r.Response(c, 500000001, nil, err.Error())
		return
	}
	arg.Default()
	arg.JwtUserMessage = r.GetUser(c)

	// 记录日志
	context := base.GetControllerBaseContext(&r.ControllerBase, c)
	context.Log.Logger.Infof("user:%+v", arg.JwtUserMessage)

	srv := services.NewPermitService(context)
	res, err := srv.AdminGroupEdit(&arg)

	if err != nil {
		r.Response(c, 500000002, nil, err.Error())
		return
	}
	r.Response(c, 0, res)
}
func (r *ControllerPermit) MenuDelete(c *gin.Context) {
	var arg pojos.ArgMenuDelete
	var err error
	err = c.Bind(&arg)
	if err != nil {
		r.Response(c, 500000001, nil, err.Error())
		return
	}
	arg.Default()
	arg.JwtUserMessage = r.GetUser(c)
	// 记录日志
	context := base.GetControllerBaseContext(&r.ControllerBase, c)
	context.Log.Logger.Infof("user:%+v", arg.JwtUserMessage)

	srv := services.NewPermitService(context)

	res, err := srv.MenuDelete(&arg)

	if err != nil {
		r.Response(c, 500000002, nil)
		return
	}
	r.Response(c, 0, res)
}

func (r *ControllerPermit) MenuSave(c *gin.Context) {
	var arg pojos.ArgMenuSave
	var err error
	err = c.ShouldBind(&arg)
	if err != nil {
		r.Response(c, 500000001, err.Error())
		return
	}
	arg.Default()
	arg.JwtUserMessage = r.GetUser(c)
	// 记录日志
	context := base.GetControllerBaseContext(&r.ControllerBase, c)

	srv := services.NewPermitService(context)
	res, err := srv.MenuSave(&arg)

	if err != nil {
		r.Response(c, 500000002, nil)
		return
	}
	r.Response(c, 0, res)
}

func (r *ControllerPermit) MenuAdd(c *gin.Context) {
	var arg pojos.ArgMenuAdd
	var err error
	err = c.ShouldBind(&arg)
	if err != nil {
		r.Response(c, 500000001, err.Error())
		return
	}
	arg.Default()
	arg.JwtUserMessage = r.GetUser(c)

	// 记录日志
	context := base.GetControllerBaseContext(&r.ControllerBase, c)

	srv := services.NewPermitService(context)
	res, err := srv.MenuAdd(&arg)

	if err != nil {
		r.Response(c, 500000002, nil)
		return
	}
	r.Response(c, 0, res)
}
func (r *ControllerPermit) AdminMenu(c *gin.Context) {
	var arg pojos.ArgAdminMenu
	err := c.Bind(&arg)
	if err != nil {
		r.Response(c, 500000001, err.Error())
		return
	}
	arg.Default()
	arg.JwtUserMessage = r.GetUser(c)

	// 记录日志
	context := base.GetControllerBaseContext(&r.ControllerBase, c)
	// context.Log.Logger.Infof("user:%+v", arg.JwtUserMessage)

	srv := services.NewPermitService(context)
	res, err := srv.AdminMenu(&arg)

	if err != nil {
		r.Response(c, 500000002, nil)
		return
	}
	r.Response(c, 0, res)
}

func (r *ControllerPermit) AdminUser(c *gin.Context) {
	var arg pojos.ArgAdminUser
	err := c.Bind(&arg)
	if err != nil {
		r.Response(c, 500000001, err.Error())
		return
	}
	arg.Default()
	arg.JwtUserMessage = r.GetUser(c)

	// 记录日志
	context := base.GetControllerBaseContext(&r.ControllerBase, c)
	context.Log.Logger.Infof("user:%+v", arg.JwtUserMessage)

	srv := services.NewPermitService(context)
	res, err := srv.AdminUser(&arg)

	if err != nil {
		r.Response(c, 500000002, nil)
		return
	}
	r.Response(c, 0, res)
}

func (r *ControllerPermit) AdminGroup(c *gin.Context) {
	var arg pojos.ArgAdminGroup
	err := c.Bind(&arg)
	if err != nil {
		r.Response(c, 500000001, nil)
		return
	}
	arg.Default()
	arg.JwtUserMessage = r.GetUser(c)

	// 记录日志
	context := base.GetControllerBaseContext(&r.ControllerBase, c)
	context.Log.Logger.Infof("user:%+v", arg.JwtUserMessage)

	srv := services.NewPermitService(context)
	res, err := srv.AdminGroup(&arg)
	if err != nil {
		r.Response(c, 500000002, nil)
		return
	}
	r.Response(c, 0, res)
}

// 获取权限菜单
func (r *ControllerPermit) Menu(c *gin.Context) {
	var arg pojos.ArgPermitMenu
	err := c.Bind(&arg)
	if err != nil {
		r.Response(c, 500000001, nil)
		return
	}
	arg.Default()
	arg.JwtUserMessage = r.GetUser(c)

	// 记录日志
	context := base.GetControllerBaseContext(&r.ControllerBase, c)
	context.Log.Logger.Infof("user:%+v", arg.JwtUserMessage)

	srv := services.NewPermitService(context)
	res, err := srv.Menu(&arg)

	if err != nil {
		r.Response(c, 500000002, nil)
		return
	}
	r.Response(c, 0, res)
}
func (r *ControllerPermit) Flag(c *gin.Context) {
	var arg pojos.ArgFlag
	err := c.Bind(&arg)
	if err != nil {
		r.Response(c, 500000000, nil)
		return
	}
	arg.JwtUserMessage = r.GetUser(c)
	srv := services.NewPermitService(base.GetControllerBaseContext(&r.ControllerBase, c))
	res, err := srv.Flag(&arg)
	if err != nil {
		r.Response(c, 500000000, nil)
		return
	}
	r.Response(c, 0, res)
}
