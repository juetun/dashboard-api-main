module github.com/juetun/dashboard-api-main

go 1.15

require (
	github.com/360EntSecGroup-Skylar/excelize v1.4.1
	github.com/aliyun/aliyun-oss-go-sdk v0.0.0-20190307165228-86c17b95fcd5
	github.com/coreos/go-systemd v0.0.0-20190719114852-fd7a80b32e1f // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/gin-gonic/gin v1.7.2
	github.com/go-errors/errors v1.4.0
	github.com/go-redis/redis/v8 v8.10.0
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/google/uuid v1.2.0
	github.com/juetun/base-wrapper v0.0.110
	github.com/lib/pq v1.3.0 // indirect
	github.com/microcosm-cc/bluemonday v1.0.1
	github.com/qiniu/go-sdk/v7 v7.9.6
	github.com/russross/blackfriday/v2 v2.0.1
	github.com/speps/go-hashids v2.0.0+incompatible
	github.com/stretchr/testify v1.7.0
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	gorm.io/gorm v1.21.11
)

replace (
	github.com/coreos/bbolt v1.3.6 => go.etcd.io/bbolt v1.3.6
	go.etcd.io/bbolt v1.3.6 => github.com/coreos/bbolt v1.3.6
)
