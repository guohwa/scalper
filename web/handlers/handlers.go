package handlers

import (
	"scalper/models"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	Handle(router *gin.Engine)
}

var handlers []Handler = []Handler{
	accountHandler,
	captchaHandler,
	customerHandler,
	homeHandler,
	passwordHandler,
	profileHandler,
	userHandler,
}

func Handle(router *gin.Engine) {
	for _, handler := range handlers {
		handler.Handle(router)
	}
}

func getUser(ctx *gin.Context) (models.User, bool) {
	u, exists := ctx.Get("user")

	if exists {
		if user, ok := u.(models.User); ok {
			return user, true
		}
	}

	return models.User{}, false
}
