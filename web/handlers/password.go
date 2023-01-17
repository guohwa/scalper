package handlers

import (
	"context"
	"net/http"

	"scalper/forms"
	"scalper/models"
	"scalper/web/handlers/response"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

var passwordHandler = &passwordhandler{}

type passwordhandler struct {
}

func (handler *passwordhandler) Handle(router *gin.Engine) {
	router.GET("/password", func(ctx *gin.Context) {
		resp := response.New(ctx)
		resp.HTML("password/index.html", response.Context{})
	})

	router.POST("/password", func(ctx *gin.Context) {
		user, ok := getUser(ctx)
		if !ok {
			ctx.Redirect(http.StatusFound, "/account/login")
			return
		}

		resp := response.New(ctx)
		form := forms.Password{}

		if err := ctx.ShouldBind(&form); err != nil {
			resp.Error(err)
			return
		}

		if user.Password != models.Encrypt(form.Password) {
			resp.Error("Invalid password")
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
			resp.Error(err)
			return
		}

		resp.Success("Password update successful", "")
	})
}
