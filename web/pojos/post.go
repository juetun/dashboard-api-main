/**
* @Author:changjiang
* @Description:
* @File:post
* @Version: 1.0.0
* @Date 2020/3/29 11:14 上午
 */
package pojos

import (
	"github.com/juetun/dashboard-api-main/web/models"
)

type PostShow struct {
	models.ZPostCate
	models.ZCategories
}
type PostTagShow struct {
	models.ZPostTag
	models.ZTags
}