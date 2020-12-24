package con_impl

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/dashboard-api-main/web/srvs/srv_impl"
)

type ControllerImg struct {
	base.ControllerBase
}

func NewControllerImg() *ControllerImg {
	controller := &ControllerImg{}
	controller.ControllerBase.Init()
	return controller
}

func (r *ControllerImg) ImgUpload(c *gin.Context) {

	file, err := c.FormFile("file")
	if err != nil {
		r.Log.Error(c, map[string]interface{}{
			"message": "post.ImgUpload",
			"err":     err.Error(),
		})
		r.Response(c, 401000004, nil)
		return
	}

	filename := filepath.Base(file.Filename)
	dst := common.ConfigUpload.ImgUploadDst + filename
	if err := c.SaveUploadedFile(file, dst); err != nil {
		r.Log.Error(c, map[string]interface{}{
			"message": "post.ImgUpload",
			"err":     err.Error(),
		})
		r.Response(c, 401000005, nil)
		return
	}

	srv := srv_impl.NewQiuNiuService(base.GetControllerBaseContext(&r.ControllerBase, c))
	// Default upload both
	data := make(map[string]interface{})
	if common.ConfigUpload.ImgUploadBoth {
		go srv.Qiniu(dst, filename)
		data["path"] = common.ConfigUpload.AppImgUrl + filename
		r.Response(c, 401000005, data)
		return
	}

	if common.ConfigUpload.QiNiuUploadImg {
		go srv.Qiniu(dst, filename)
		data["path"] = common.ConfigUpload.QiNiuHostName + filename
		r.Response(c, 401000005, data)
		return
	}

	data["path"] = common.ConfigUpload.AppImgUrl + filename
	r.Response(c, 401000005, data)
	return
}
