/**
* @Author:changjiang
* @Description:
* @File:file_upload
* @Version: 1.0.0
* @Date 2020/6/11 11:33 下午
 */
package export

import (
	"path/filepath"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/juetun/dashboard-api-main/basic/plugins"
)

type FileUpload struct {
	FilePath    string `json:"file_path"` // <yourLocalFileName>由本地文件路径加文件名包括后缀组成，例如/users/local/myfile.txt。
	DownloadUrl string `json:"download_url"`
	Endpoint    string `json:"endpoint"`
	ObjectName  string `json:"object_name"`
	Err         error
}

func NewNewFileUpload() (r *FileUpload) {
	r = &FileUpload{}
	return
}

func (r *FileUpload) Run() {
	var data = plugins.GetOssConfig()

	// // Endpoint以杭州为例，其它Region请按实际情况填写。
	// endpoint := "http://oss-cn-hangzhou.aliyuncs.com"
	// // 阿里云主账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM账号进行API访问或日常运维，请登录 https://ram.console.aliyun.com 创建RAM账号。
	// accessKeyId := "<yourAccessKeyId>"
	// accessKeySecret := "<yourAccessKeySecret>"
	// bucketName := "<yourBucketName>"
	// // <yourObjectName>上传文件到OSS时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
	// objectName := "<yourObjectName>"

	var client *oss.Client
	var bucket *oss.Bucket

	r.Endpoint = data.Endpoint
	// 创建OSSClient实例。
	client, r.Err = oss.New(r.Endpoint, data.AccessKeyId, data.AccessKeySecret)
	if r.Err != nil {
		return
	}

	// 获取存储空间。
	bucket, r.Err = client.Bucket(data.BucketName)
	if r.Err != nil {
		return
	}
	rulePrefix := data.DirName + "/export/"
	r.ObjectName = rulePrefix + strings.TrimPrefix(r.FilePath, filepath.Dir(r.FilePath))

	// 设置导出 文件14天过期
	rule1 := oss.BuildLifecycleRuleByDays("rule1", rulePrefix, true, 14)
	rules := []oss.LifecycleRule{rule1}
	r.Err = client.SetBucketLifecycle(data.BucketName, rules)
	if r.Err != nil {
		return
	}

	// 上传文件。
	r.Err = bucket.PutObjectFromFile(r.ObjectName, r.FilePath)
	if r.Err != nil {
		return
	}
	r.DownloadUrl = r.Endpoint + "/" + r.ObjectName
}

func (r *FileUpload) SetFile(fileName string) *FileUpload {
	r.FilePath = fileName
	return r
}
