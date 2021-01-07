package srv_impl

import (
	"context"

	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/base-wrapper/lib/utils"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

// 自定义返回值结构体
type MyPutRet struct {
	Key    string
	Hash   string
	Fsize  int
	Bucket string
	Name   string
}

type QiuNiuService struct {
	base.ServiceBase
}

func NewQiuNiuService(context ...*base.Context) (p *QiuNiuService) {
	p = &QiuNiuService{}
	p.SetContext(context...)
	return
}

// Upload file to Qiniu
// LocalFile is the local file, such as "./static/images/uploads/2.jpeg"
// FileName is the name what  qiniu name as
// The storage Zone is default
func (r *QiuNiuService) Qiniu(localFile string, fileName string) {
	accessKey := common.ConfigUpload.QiNiuAccessKey
	secretKey := common.ConfigUpload.QiNiuSecretKey
	// localFile := "./static/images/uploads/2.jpeg"
	bucket := common.ConfigUpload.QiNiuBucket
	key := fileName
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	switch common.ConfigUpload.QiNiuZone {
	case "HUABEI":
		cfg.Zone = &storage.ZoneHuabei
	case "HUADONG":
		cfg.Zone = &storage.ZoneHuadong
	case "BEIMEI":
		cfg.Zone = &storage.ZoneBeimei
	case "HUANAN":
		cfg.Zone = &storage.ZoneHuanan
	case "XINJIAPO":
		cfg.Zone = &storage.ZoneXinjiapo
	default:
		cfg.Zone = &storage.ZoneHuabei
	}
	cfg.UseHTTPS = false
	cfg.UseCdnDomains = false

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	putExtra := storage.PutExtra{
		Params: map[string]string{
			// "x:name": "github logo",
		},
	}
	err := formUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)
	if err != nil {
		r.Context.Error(
			map[string]interface{}{"message": "service.QiNiu upload file",
				"err": err.Error(),
			},
		)
		utils.Alarm("文件上传七牛失败了,文件名是" + fileName)
		return
	}
	// fmt.Println("234222",ret,"最后是谁啥")
	// fmt.Println(ret.Key,"234222", ret.Hash,"最后是谁啥")
	r.Context.Info(
		map[string]interface{}{"message": "service.QiNiu upload file",
			"end": "succeed",
		},
	)

	return
}
