package handlers

import (
	"context"

	"scalper/forms"
	"scalper/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var accountHandler = &accounthandler{}

type accounthandler struct {
	base
}

func (handler *accounthandler) Handle(router *gin.Engine) {
	router.GET("/account/login", func(ctx *gin.Context) {
		handler.HTML(ctx, "account/login.html", nil)
	})

	router.POST("/account/login", func(ctx *gin.Context) {
		form := forms.Login{}
		if err := ctx.ShouldBind(&form); err != nil {
			handler.Error(ctx, err)
			return
		}

		if !captchaHandler.Verify(ctx, form.Verify) {
			handler.Error(ctx, "Invalid verify code")
			return
		}

		user := models.User{}
		if err := models.UserCollection.FindOne(context.Background(), bson.M{"username": form.Username}).Decode(&user); err != nil {
			if err == mongo.ErrNoDocuments {
				handler.Error(ctx, "User Does not exists.")
			} else {
				handler.Error(ctx, err)
			}
			return
		}

		if user.Password != models.Encrypt(form.Password) {
			handler.Error(ctx, "Invalid password")
			return
		}

		if user.Status == "Disable" {
			handler.Error(ctx, "User disabled")
			return
		}

		session := sessions.Default(ctx)
		session.Set("user-id", user.ID.Hex())

		if err := session.Save(); err == nil {
			handler.Success(ctx, "Login Successful", "/")
		} else {
			handler.Error(ctx, "Login Failed")
		}
	})

	router.GET("/account/logout", func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		session.Delete("user-id")

		if err := session.Save(); err == nil {
			handler.Success(ctx, "Logout successful", "/")
		} else {
			handler.Error(ctx, "Logout failed")
		}
	})

}
