package controllers

import "github.com/gin-gonic/gin"

type ConsoleAuth interface {
	Register(*gin.Context)
	AuthRegister(*gin.Context)
	Login(*gin.Context)
	AuthLogin(*gin.Context)
	Logout(*gin.Context)
	DelCache(*gin.Context)
}
