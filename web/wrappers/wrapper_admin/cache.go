package wrapper_admin

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/base-wrapper/lib/utils"
	"github.com/juetun/dashboard-api-main/web/models"
)

type (
	ArgReloadAppCacheConfig struct {

	}
	ResultReloadAppCacheConfig struct {
		Result bool `json:"result"`
	}
	ArgCacheParamList struct {
		HeaderInfo common.HeaderInfo
		response.PageQuery
		MicroApp string `json:"micro_app" form:"micro_app"`
	}
	ResultCacheParamList struct {
		Pager response.Pager `json:"pager"`
	}
	ResultCacheParamItem struct {
		Id        int64
		Key       string
		Expire    string
		MicroApp  string
		Desc      string
		CreatedAt string
		UpdatedAt string
	}
	ArgClearCache struct {
	}
	ResultClearCache struct {
		Result bool `json:"result"`
	}
)

func (r *ArgReloadAppCacheConfig) Default(ctx *base.Context) (err error) {
	return
}

func (r *ResultCacheParamItem) ParseWithCacheParamData(data *models.CacheKeyData) {
	r.Id = data.Id
	r.Key = data.Key
	r.Expire = data.Expire
	r.MicroApp = data.MicroApp
	r.Desc = data.Desc
	r.CreatedAt = data.CreatedAt.Format(utils.DateTimeGeneral)
	r.UpdatedAt = data.UpdatedAt.Format(utils.DateTimeGeneral)
	return
}

func (r *ArgCacheParamList) Default(ctx *base.Context) (err error) {

	return
}

func (r *ArgClearCache) Default(ctx *base.Context) (err error) {

	return
}
