// Package admin_impl
package admin_impl

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/cons/admins"
	"github.com/juetun/dashboard-api-main/web/srvs/srv_impl"
	"github.com/juetun/dashboard-api-main/web/wrappers"
	"github.com/juetun/dashboard-api-main/web/wrappers/wrapper_admin"
)

type ControllerPermit struct {
	base.ControllerBase
}

func NewControllerPermit() admins.Permit {
	controller := &ControllerPermit{}
	controller.ControllerBase.Init()
	return controller
}

func (r *ControllerPermit) MenuImportSet(c *gin.Context) {
	var (
		res *wrappers.ResultMenuImportSet
		arg wrappers.ArgMenuImportSet
	)
	var err error

	if err = c.Bind(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	if err = arg.Default(c); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	// 记录日志
	if res, err = srv_impl.NewSrvPermitGroupImpl(base.CreateContext(&r.ControllerBase, c)).
		MenuImportSet(&arg); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)

}

func (r *ControllerPermit) GetMenu(c *gin.Context) {
	var (
		err error
		arg wrappers.ArgGetMenu
		res wrappers.ResultGetMenu
	)

	if err = c.Bind(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	if err = arg.Default(c); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}

	// 记录日志
	if res, err = srv_impl.NewPermitServiceImpl(base.CreateContext(&r.ControllerBase, c)).
		GetMenu(&arg); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)

}
func (r *ControllerPermit) AdminMenuSearch(c *gin.Context) {
	var (
		arg wrappers.ArgAdminMenu
		err error
		res wrappers.ResAdminMenuSearch
	)

	if err = c.Bind(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	if err = arg.Default(c); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}

	// 记录日志
	if res, err = srv_impl.NewPermitServiceImpl(base.CreateContext(&r.ControllerBase, c)).
		AdminMenuSearch(&arg); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)

}

func (r *ControllerPermit) AdminUserGroupRelease(c *gin.Context) {
	var arg wrappers.ArgAdminUserGroupRelease
	var err error
	err = c.Bind(&arg)
	if err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	if err = arg.Default(c); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}

	// 记录日志
	srv := srv_impl.NewSrvPermitGroupImpl(base.CreateContext(&r.ControllerBase, c))
	res, err := srv.AdminUserGroupRelease(&arg)

	if err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)
}

func (r *ControllerPermit) AdminUserGroupAdd(c *gin.Context) {
	var arg wrappers.ArgAdminUserGroupAdd
	var err error

	if err = c.Bind(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	if err = arg.Default(c); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)

		return
	}
	res, err := srv_impl.NewSrvPermitGroupImpl(base.CreateContext(&r.ControllerBase, c)).
		AdminUserGroupAdd(&arg)

	if err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	r.Response(c, base.SuccessCode, res)
}

func (r *ControllerPermit) AdminGroupDelete(c *gin.Context) {
	var arg wrappers.ArgAdminGroupDelete
	var err error
	err = c.Bind(&arg)
	if err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	if err = arg.Default(c); err != nil {
		return
	}

	// 记录日志
	context := base.CreateContext(&r.ControllerBase, c)
	context.Info(map[string]interface{}{"arg": arg})

	res, err := srv_impl.NewSrvPermitGroupImpl(context).
		AdminGroupDelete(&arg)

	if err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)
}

func (r *ControllerPermit) AdminGroupEdit(c *gin.Context) {
	var arg wrappers.ArgAdminGroupEdit
	var err error
	err = c.Bind(&arg)
	if err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	if err = arg.Default(c); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}

	// 记录日志
	res, err := srv_impl.NewSrvPermitGroupImpl(base.CreateContext(&r.ControllerBase, c)).
		AdminGroupEdit(&arg)
	if err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)
}

func (r *ControllerPermit) DeleteImport(c *gin.Context) {
	var err error
	var arg wrappers.ArgDeleteImport
	if err = c.ShouldBindUri(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	var res *wrappers.ResultDeleteImport
	// 记录日志
	if res, err = srv_impl.NewSrvPermitImport(base.CreateContext(&r.ControllerBase, c)).
		DeleteImport(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	} else {
		r.Response(c, base.SuccessCode, res)
	}
}

func (r *ControllerPermit) UpdateImportValue(c *gin.Context) {
	var arg wrappers.ArgUpdateImportValue
	var err error
	var res *wrappers.ResultUpdateImportValue
	if err = c.Bind(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	if err = arg.Default(c); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}

	// 记录日志
	if res, err = srv_impl.NewSrvPermitImport(base.CreateContext(&r.ControllerBase, c)).
		UpdateImportValue(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	r.Response(c, base.SuccessCode, res)
	return
}
func (r *ControllerPermit) ImportList(c *gin.Context) {
	var arg wrappers.ArgImportList
	var err error
	var res *wrappers.ResultImportList
	if err = c.ShouldBind(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	if err = arg.Default(c); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}

	// 记录日志
	if res, err = srv_impl.NewSrvPermitImport(base.CreateContext(&r.ControllerBase, c)).
		ImportList(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	r.Response(c, base.SuccessCode, res)
	return
}
func (r *ControllerPermit) EditImport(c *gin.Context) {
	var arg wrappers.ArgEditImport
	var err error
	var res *wrappers.ResultEditImport
	if err = c.Bind(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	if err = arg.Default(c); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}

	// 记录日志
	if res, err = srv_impl.NewSrvPermitImport(base.CreateContext(&r.ControllerBase, c)).
		EditImport(&arg); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)
	return
}
func (r *ControllerPermit) GetImport(c *gin.Context) {
	var (
		res *wrappers.ResultGetImport
		arg wrappers.ArgGetImport
		err error
	)

	if err = c.ShouldBind(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	if err = arg.Default(c); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}

	// 记录日志
	if res, err = srv_impl.NewSrvPermitImport(base.CreateContext(&r.ControllerBase, c)).
		GetImport(&arg); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)

}

// GetImportByMenuId 根据菜单号 获取页面的接口ID
//
func (r *ControllerPermit) GetImportByMenuId(c *gin.Context) {
	var (
		arg wrappers.ArgGetImportByMenuId
		res = wrappers.ResultGetImportByMenuId{}
		err error
	)

	if err = c.ShouldBind(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	if err = arg.Default(c); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}

	// 记录日志
	if res, err = srv_impl.NewSrvPermitImport(base.CreateContext(&r.ControllerBase, c)).
		GetImportByMenuId(&arg); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)

}

func (r *ControllerPermit) MenuImport(c *gin.Context) {
	var (
		arg wrapper_admin.ArgMenuImport
		err error
		res *wrapper_admin.ResultMenuImport
	)

	if err = c.ShouldBind(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	if err = arg.Default(c); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	// 记录日志
	if res, err = srv_impl.NewPermitServiceImpl(base.CreateContext(&r.ControllerBase, c)).
		MenuImport(&arg); err != nil {
		r.ResponseError(c, err)
		return
	} else {
		r.Response(c, base.SuccessCode, res)
	}
}
func (r *ControllerPermit) MenuDelete(c *gin.Context) {

	var (
		arg wrappers.ArgMenuDelete
		err error
		res *wrappers.ResultMenuDelete
	)

	if err = c.Bind(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	if err = arg.Default(c); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}

	// 记录日志
	if res, err = srv_impl.NewPermitServiceImpl(base.CreateContext(&r.ControllerBase, c)).
		MenuDelete(&arg); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)

}

func (r *ControllerPermit) MenuSave(c *gin.Context) {
	var arg wrappers.ArgMenuSave
	var err error
	if err = c.ShouldBind(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	if err = arg.Default(c); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}

	var res *wrappers.ResultMenuSave
	// 记录日志
	if res, err = srv_impl.NewPermitServiceImpl(base.CreateContext(&r.ControllerBase, c)).
		MenuSave(&arg); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)
}

func (r *ControllerPermit) MenuAdd(c *gin.Context) {
	var arg wrappers.ArgMenuAdd
	var err error
	if err = c.ShouldBind(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}

	if err = arg.Default(c); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)

		return
	}

	// 记录日志
	res, err := srv_impl.NewPermitServiceImpl(base.CreateContext(&r.ControllerBase, c)).
		MenuAdd(&arg)
	if err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)
}

func (r *ControllerPermit) AdminMenuWithCheck(c *gin.Context) {
	var (
		arg wrappers.ArgAdminMenuWithCheck
		res *wrappers.ResultMenuWithCheck
		err error
	)
	if err = c.Bind(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	if err = arg.Default(c); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}

	if res, err = srv_impl.NewSrvPermitMenu(base.CreateContext(&r.ControllerBase, c)).
		AdminMenuWithCheck(&arg); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)
}

func (r *ControllerPermit) AdminMenu(c *gin.Context) {
	var arg wrappers.ArgAdminMenu
	err := c.Bind(&arg)
	if err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	if err = arg.Default(c); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)

		return
	}

	// 记录日志
	res, err := srv_impl.NewSrvPermitMenu(base.CreateContext(&r.ControllerBase, c)).
		AdminMenu(&arg)

	if err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)
}

// AdminSetPermit 设置权限
func (r *ControllerPermit) AdminSetPermit(c *gin.Context) {
	var arg wrappers.ArgAdminSetPermit
	var err error
	var res *wrappers.ResultAdminSetPermit

	if err = c.Bind(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	if err = arg.Default(c); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}

	// 记录日志
	if res, err = srv_impl.NewPermitServiceImpl(base.CreateContext(&r.ControllerBase, c)).
		AdminSetPermit(&arg); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)
}

func (r *ControllerPermit) AdminGroup(c *gin.Context) {
	var arg wrappers.ArgAdminGroup
	var err error
	if err = c.Bind(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	if err = arg.Default(c); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}

	// 记录日志
	res, err := srv_impl.NewSrvPermitGroupImpl(base.CreateContext(&r.ControllerBase, c)).
		AdminGroup(&arg)
	if err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)
}

// Menu 获取权限菜单
func (r *ControllerPermit) Menu(c *gin.Context) {
	var (
		err error
		arg wrappers.ArgPermitMenu
		res *wrappers.ResultPermitMenuReturn
	)
	if err = c.Bind(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	if err = arg.Default(c); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}

	// 记录日志
	if res, err = srv_impl.NewSrvPermitMenu(base.CreateContext(&r.ControllerBase, c)).
		Menu(&arg); err != nil {
		r.ResponseError(c, err)
		return
	}
	r.Response(c, base.SuccessCode, res)
}

func (r *ControllerPermit) GetAppConfig(c *gin.Context) {

	var (
		arg wrappers.ArgGetAppConfig
		res *wrappers.ResultGetAppConfig
		err error
	)

	if err = c.ShouldBind(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	if err = arg.Default(c); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}

	if res, err = srv_impl.NewSrvPermitAppImpl(base.CreateContext(&r.ControllerBase, c)).
		GetAppConfig(&arg); err != nil {
		r.ResponseError(c, err, base.ErrorParameterCode)
		return
	}
	r.Response(c, base.SuccessCode, res)

}
