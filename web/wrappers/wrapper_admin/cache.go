package wrapper_admin

import (
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/base-wrapper/lib/common/response"
	"github.com/juetun/base-wrapper/lib/utils"
	"github.com/juetun/dashboard-api-main/web/models"
	"github.com/juetun/library/common/app_param"
)

type (
	ArgReloadAppCacheConfig struct {
		app_param.ArgGetCacheParamConfig
		TimeNow base.TimeNormal `json:"-" form:"-"`
	}
	ResultReloadAppCacheConfig struct {
		Result bool `json:"result"`
	}
	ArgCacheParamList struct {
		HeaderInfo common.HeaderInfo
		response.PageQuery
		MicroApp string `json:"micro_app" form:"micro_app"`
		Key      string `json:"key" form:"key"`
	}
	ResultCacheParamList struct {
		Pager response.Pager `json:"pager"`
	}
	ResultCacheParamItem struct {
		Id                int64  `json:"id"`
		Key               string `json:"key"`
		Expire            string `json:"expire"`
		MicroApp          string `json:"micro_app"`
		Desc              string `json:"desc"`
		CacheDataType     uint8  `json:"cache_data_type"`
		CacheDataTypeName string `json:"cache_data_type_name"`
		CreatedAt         string `json:"created_at"`
		UpdatedAt         string `json:"updated_at"`
	}
	ArgClearCache struct {
	}
	ResultClearCache struct {
		Result bool `json:"result"`
	}
)

func (r *ArgReloadAppCacheConfig) Default(ctx *base.Context) (err error) {
	r.TimeNow = base.GetNowTimeNormal()
	return
}

func (r *ResultCacheParamItem) ParseWithCacheParamData(data *models.CacheKeyData) {
	r.Id = data.Id
	r.Key = data.Key
	r.Expire = data.Expire
	r.MicroApp = data.MicroApp
	r.Desc = data.Desc
	r.CacheDataType = data.CacheDataType
	r.CacheDataTypeName = data.ParseCacheDataType()
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
