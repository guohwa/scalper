package handlers

import (
	"context"

	"scalper/config"
	"scalper/forms"
	"scalper/log"
	"scalper/models"
	"scalper/web/handlers/response"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

		if form.Mode != gin.Mode() {
			gin.SetMode(form.Mode)
		}

		if form.Level != log.GetLevel().String() {
			level, err := logrus.ParseLevel(form.Level)
			if err != nil {
				resp.Error(err)
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
			resp.Error(err)
			return
		}

		resp.Success("App update successful", "")
	})
}
