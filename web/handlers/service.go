package handlers

import (
	"scalper/service"

	"github.com/gin-gonic/gin"
)

var serviceHandler = &servicehandler{}

type servicehandler struct {
	base
}

func (handler *servicehandler) Handle(router *gin.Engine) {
	router.GET("/service", func(ctx *gin.Context) {
		handler.HTML(ctx, "service/index.html", Context{
			"status": service.Status(),
		})
	})

	router.GET("/service/start", func(ctx *gin.Context) {
		if err := service.Start(); err != nil {
			handler.Error(ctx, err)
			return
		}

		handler.Success(ctx, "service start successful", "")
	})

	router.GET("/service/stop", func(ctx *gin.Context) {
		service.Stop()

		handler.Success(ctx, "service stop successful", "")
	})

}
