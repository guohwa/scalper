package handlers

import (
	"context"
	"scalper/config"
	"scalper/forms"
	"scalper/models"
	"scalper/web/handlers/response"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

var appHandler = &apphandler{}

type apphandler struct {
}

func (handler *apphandler) Handle(router *gin.Engine) {
	router.GET("/app", func(ctx *gin.Context) {
		resp := response.New(ctx)
		resp.HTML("app/index.html", response.Context{
			"app": config.App,
		})
	})

	router.POST("/app", func(ctx *gin.Context) {
		resp := response.New(ctx)
		form := forms.App{}

		if err := ctx.ShouldBind(&form); err != nil {
			resp.Error(err)
			return
		}

		update := bson.M{"$set": bson.M{
			"title": form.Title,
			"mode":  form.Mode,
		}}
		if _, err := models.AppCollection.UpdateByID(
			context.TODO(),
			config.App.ID,
			update,
		); err != nil {
			resp.Error(err)
			return
		}

		config.App.Title = form.Title
		config.App.Mode = form.Mode

		if form.Mode != gin.Mode() {
			gin.SetMode(form.Mode)
		}

		resp.Success("App update successful", "")
	})
}
