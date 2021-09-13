/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2018-12-27
 * Time: 00:07
 */
package admins

import "github.com/gin-gonic/gin"

type Console interface {
	//
	Index(*gin.Context)
	//
	Create(*gin.Context)

	//
	Store(*gin.Context)

	//
	Edit(*gin.Context)

	//
	Update(*gin.Context)

	//
	Destroy(*gin.Context)
}
