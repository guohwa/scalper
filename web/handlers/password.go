package handlers

import (
	"context"
	"net/http"

	"scalper/forms"
	"scalper/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

var passwordHandler = &passwordhandler{}

type passwordhandler struct {
	base
}

func (handler *passwordhandler) Handle(router *gin.Engine) {
	router.GET("/password", func(ctx *gin.Context) {
		handler.HTML(ctx, "password/index.html", Context{})
	})

	router.POST("/password", func(ctx *gin.Context) {
		user, ok := getUser(ctx)
		if !ok {
			ctx.Redirect(http.StatusFound, "/account/login")
			return
		}

		form := forms.Password{}

		if err := ctx.ShouldBind(&form); err != nil {
			handler.Error(ctx, err)
			return
		}

		if user.Password != models.Encrypt(form.Password) {
			handler.Error(ctx, "Invalid password")
			return
		}

		update := bson.M{"$set": bson.M{
			"password": models.Encrypt(form.NewPassword),
		}}
		if _, err := models.UserCollection.UpdateByID(
			context.TODO(),
			user.ID,
			update,
		); err != nil {
			handler.Error(ctx, err)
			return
		}

		handler.Success(ctx, "Password update successful", "")
	})
}
