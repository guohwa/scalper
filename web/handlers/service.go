package handlers

import (
	"scalper/web/handlers/response"

	"github.com/gin-gonic/gin"
)

var serviceHandler = &servicehandler{}

type servicehandler struct {
}

func (handler *servicehandler) Handle(router *gin.Engine) {
	router.GET("/service", func(ctx *gin.Context) {
		resp := response.New(ctx)

		resp.HTML("service/index.html", response.Context{})
	})
}
