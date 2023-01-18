package handlers

import (
	"scalper/service"
	"scalper/web/handlers/response"

	"github.com/gin-gonic/gin"
)

var serviceHandler = &servicehandler{}

type servicehandler struct {
}

func (handler *servicehandler) Handle(router *gin.Engine) {
	router.GET("/service", func(ctx *gin.Context) {
		resp := response.New(ctx)

		resp.HTML("service/index.html", response.Context{
			"status": service.Status(),
		})
	})

	router.GET("/service/start", func(ctx *gin.Context) {
		resp := response.New(ctx)

		if err := service.Start(); err != nil {
			resp.Error(err)
			return
		}

		resp.Success("service start successful", "")
	})

	router.GET("/service/stop", func(ctx *gin.Context) {
		resp := response.New(ctx)

		service.Stop()

		resp.Success("service stop successful", "")
	})

}
