package handlers

import (
	"scalper/web/handlers/response"

	"github.com/gin-gonic/gin"
)

var profileHandler = &profilehandler{}

type profilehandler struct {
}

func (handler *profilehandler) Handle(router *gin.Engine) {
	router.GET("/profile", func(ctx *gin.Context) {
		resp := response.New(ctx)
		resp.HTML("profile/index.html", response.Context{})
	})
}
