package handlers

import (
	"github.com/gin-gonic/gin"
)

var profileHandler = &profilehandler{}

type profilehandler struct {
	base
}

func (handler *profilehandler) Handle(router *gin.Engine) {
	router.GET("/profile", func(ctx *gin.Context) {
		handler.HTML(ctx, "profile/index.html", Context{})
	})
}
