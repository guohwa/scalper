package handlers

import (
	"context"

	"scalper/config"
	"scalper/forms"
	"scalper/log"
	"scalper/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

var appHandler = &apphandler{}

type apphandler struct {
	base
}

func (handler *apphandler) Handle(router *gin.Engine) {
	router.GET("/app", func(ctx *gin.Context) {
		handler.HTML(ctx, "app/index.html", Context{
			"app": config.App,
		})
	})

	router.POST("/app", func(ctx *gin.Context) {
		form := forms.App{}

		if err := ctx.ShouldBind(&form); err != nil {
			handler.Error(ctx, err)
			return
		}

		if form.Mode != gin.Mode() {
			gin.SetMode(form.Mode)
		}

		if form.Level != log.GetLevel().String() {
			level, err := log.ParseLevel(form.Level)
			if err != nil {
				handler.Error(ctx, err)
				return
			}
			log.SetLevel(level)
		}

		config.App.Title = form.Title
		config.App.Mode = form.Mode
		config.App.Level = form.Level

		update := bson.M{"$set": bson.M{
			"title": form.Title,
			"mode":  form.Mode,
			"level": form.Level,
		}}
		if _, err := models.AppCollection.UpdateByID(
			context.TODO(),
			config.App.ID,
			update,
		); err != nil {
			handler.Error(ctx, err)
			return
		}

		handler.Success(ctx, "App update successful", "")
	})
}
