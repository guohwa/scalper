package middleware

import (
	"scalper/config"
	"time"

	"github.com/gin-gonic/gin"
)

func Global() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("title", config.App.Title)
		ctx.Set("now", time.Now())
	}
}
