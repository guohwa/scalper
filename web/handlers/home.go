package handlers

import (
	"scalper/web/handlers/response"

	"github.com/gin-gonic/gin"
)

var homeHandler = &homehandler{}

type homehandler struct {
}

func (handler *homehandler) Handle(router *gin.Engine) {
	router.GET("/", func(ctx *gin.Context) {
		ctx.Request.URL.Path = "/home"
		router.HandleContext(ctx)
	})

	router.GET("/home", func(ctx *gin.Context) {
		resp := response.New(ctx)
		resp.HTML("home/index.html", response.Context{})
	})
}
