module github.com/juetun/dashboard-api-main

go 1.15

require (
	github.com/360EntSecGroup-Skylar/excelize v1.4.1
	github.com/aliyun/aliyun-oss-go-sdk v0.0.0-20190307165228-86c17b95fcd5
	github.com/gin-gonic/gin v1.8.1
	github.com/go-errors/errors v1.4.0
	github.com/go-redis/redis/v8 v8.10.0
	github.com/juetun/base-wrapper v0.0.216
	github.com/microcosm-cc/bluemonday v1.0.1
	github.com/qiniu/go-sdk/v7 v7.9.6
	github.com/russross/blackfriday/v2 v2.0.1
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	gorm.io/driver/mysql v1.3.4 // indirect
	gorm.io/gorm v1.23.4
)

replace (
	github.com/coreos/bbolt v1.3.6 => go.etcd.io/bbolt v1.3.6
	go.etcd.io/bbolt v1.3.6 => github.com/coreos/bbolt v1.3.6
)
