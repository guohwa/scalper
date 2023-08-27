package handlers

import (
	"github.com/gin-gonic/gin"
)

var homeHandler = &homehandler{}

type homehandler struct {
	base
}

func (handler *homehandler) Handle(router *gin.Engine) {
	router.GET("/", func(ctx *gin.Context) {
		ctx.Request.URL.Path = "/home"
		router.HandleContext(ctx)
	})

	router.GET("/home", func(ctx *gin.Context) {
		handler.HTML(ctx, "home/index.html", Context{})
	})
}
