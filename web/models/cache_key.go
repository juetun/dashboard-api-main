package models

import (
	"encoding/json"
	"fmt"
	"github.com/juetun/base-wrapper/lib/base"
)

type (
	CacheKeyData struct {
		Id        int64            `gorm:"column:id;primary_key" json:"id"`
		Key       string           `gorm:"column:key;type:varchar(150);not null;default:'';comment:缓存的key" json:"key"`            // key
		Expire    string           `gorm:"column:expire;type:varchar(50);not null;default:'';comment:有效期" json:"expire"`          // 过期时间
		MicroApp  string           `gorm:"column:micro_app;type:varchar(100);not null;default:'';comment:微服key" json:"micro_app"` // 服务
		Desc      string           `gorm:"column:desc;type:varchar(500);not null;default:'';comment:缓存描述" json:"desc"`            // 缓存使用场景描述
		CreatedAt base.TimeNormal  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
		UpdatedAt base.TimeNormal  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
		DeletedAt *base.TimeNormal `gorm:"column:deleted_at;" json:"-"`
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

//实现 序列化方法 encoding.BinaryMarshaler
func (r *CacheKeyData) MarshalBinary() (data []byte, err error) {
	if r == nil {
		data = []byte{}
		return
	}
	data, err = json.Marshal(r)
	return
}
