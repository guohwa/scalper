package handlers

import (
	"scalper/web/handlers/response"

	"github.com/gin-gonic/gin"
)

var homeHandler = &homehandler{}

type homehandler struct {
}

func (handler *homehandler) Handle(router *gin.Engine) {
	router.GET("/", handler.render)

	router.GET("/home", handler.render)
}

func (handler *homehandler) render(ctx *gin.Context) {
	resp := response.New(ctx)
	resp.HTML("home/index.html", response.Context{})
}
