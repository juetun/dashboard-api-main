package models

import (
	"encoding/json"
	"fmt"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/library/common/app_param"
)

type (
	CacheKeyData struct {
		Id            int64            `gorm:"column:id;primary_key" json:"id"`
		Key           string           `gorm:"column:key;type:varchar(150);not null;default:'';uniqueIndex:idx_id,priority:2;comment:缓存的key" json:"key"`            // key
		Expire        string           `gorm:"column:expire;type:varchar(50);not null;default:'';comment:有效期" json:"expire"`                                        // 过期时间
		MicroApp      string           `gorm:"column:micro_app;type:varchar(100);not null;default:'';uniqueIndex:idx_id,priority:1;comment:微服key" json:"micro_app"` // 服务
		Desc          string           `gorm:"column:desc;type:varchar(500);not null;default:'';comment:缓存描述" json:"desc"`                                          // 缓存使用场景描述
		CacheDataType uint8            `gorm:"column:cache_data_type;type:tinyint(2);not null;default:1;comment:缓存数据类型" json:"cache_data_type"`
		CreatedAt     base.TimeNormal  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
		UpdatedAt     base.TimeNormal  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
		DeletedAt     *base.TimeNormal `gorm:"column:deleted_at;" json:"-"`
	}
)

func (r *CacheKeyData) TableName() string {
	return fmt.Sprintf("%scache_key_data", TablePrefix)
}

func (r *CacheKeyData) GetTableComment() (res string) {
	res = "缓存信息Key记录表"
	return
}

func (r *CacheKeyData) UnmarshalBinary(data []byte) (err error) {
	if len(data) == 0 {
		r = &CacheKeyData{}
	}
	err = json.Unmarshal(data, r)
	return
}

func (r *CacheKeyData) ParseCacheDataType() (res string) {

	sliceCacheDataType, _ := app_param.SliceCacheDataType.GetMapAsKeyUint8()
	if dt, ok := sliceCacheDataType[r.CacheDataType]; ok {
		res = dt
		return
	}
	res = fmt.Sprintf("未知(%d)", r.CacheDataType)
	return
}

func (r *CacheKeyData) Default() (err error) {

	if r.CacheDataType == 0 {
		r.CacheDataType = app_param.CacheDataTypeHashString
	}
	return
}

//实现 序列化方法 encoding.BinaryMarshaler
func (r *CacheKeyData) MarshalBinary() (data []byte, err error) {
	if r == nil {
		data = []byte{}
		return
	}
	data, err = json.Marshal(r)
	return
}
