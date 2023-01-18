package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	forms "scalper/forms/customer"
	"scalper/models"
	"scalper/web/handlers/response"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var customerHandler = &customerhandler{}

type customerhandler struct {
}

func (handler *customerhandler) Handle(router *gin.Engine) {
	router.GET("/customer", func(ctx *gin.Context) {
		user, ok := getUser(ctx)
		if !ok {
			ctx.Redirect(http.StatusFound, "/account/login")
			return
		}

		resp := response.New(ctx)
		filter := bson.M{
			"userId": user.ID,
		}

		count, err := models.CustomerCollection.CountDocuments(
			context.TODO(),
			filter,
			options.Count().SetMaxTime(2*time.Second))
		if err != nil {
			resp.Error(err)
			return
		}

		page, err := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
		if err != nil {
			resp.Error(err)
			return
		}
		limit, err := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)
		if err != nil {
			resp.Error(err)
			return
		}
		cursor, err := models.CustomerCollection.Find(
			context.TODO(),
			filter, options.Find().SetSkip((page-1)*limit).SetLimit(limit),
		)
		if err != nil {
			resp.Error(err)
			return
		}

		var items []models.Customer
		if err = cursor.All(context.TODO(), &items); err != nil {
			resp.Error(err)
			return
		}

		resp.HTML("customer/index.html", response.Context{
			"page":  page,
			"count": count,
			"limit": limit,
			"items": items,
		})
	})

	router.GET("/customer/add", func(ctx *gin.Context) {
		resp := response.New(ctx)
		resp.HTML("customer/add.html", response.Context{})
	})

	router.POST("/customer/save", func(ctx *gin.Context) {
		user, ok := getUser(ctx)
		if !ok {
			ctx.Redirect(http.StatusFound, "/account/login")
			return
		}

		resp := response.New(ctx)
		form := forms.Save{}

		if err := ctx.ShouldBind(&form); err != nil {
			resp.Error(err)
			return
		}

		saved := bson.M{
			"userId":    user.ID,
			"name":      form.Name,
			"apiKey":    form.ApiKey,
			"apiSecret": form.ApiSecret,
			"capital":   form.Capital,
			"position":  "",
			"status":    form.Status,
		}
		if _, err := models.CustomerCollection.InsertOne(
			context.TODO(),
			saved,
		); err != nil {
			resp.Error(err)
			return
		}

		resp.Success("Customer save successful", "/customer")
	})

	router.GET("/customer/edit/:id", func(ctx *gin.Context) {
		resp := response.New(ctx)

		sId := ctx.Param("id")
		uId, err := primitive.ObjectIDFromHex(sId)
		if err != nil {
			resp.Error(err)
			return
		}

		customer := models.Customer{}
		if err := models.CustomerCollection.FindOne(context.TODO(), bson.M{
			"_id": uId,
		}).Decode(&customer); err != nil {
			resp.Error(err)
			return
		}

		resp.HTML("customer/edit.html", response.Context{
			"item": customer,
		})
	})

	router.POST("/customer/update/:id", func(ctx *gin.Context) {
		resp := response.New(ctx)
		form := forms.Update{}
		sId := ctx.Param("id")
		uId, err := primitive.ObjectIDFromHex(sId)
		if err != nil {
			resp.Error(err)
			return
		}

		if err := ctx.ShouldBind(&form); err != nil {
			resp.Error(err)
			return
		}

		update := bson.M{"$set": bson.M{
			"name":      form.Name,
			"apiKey":    form.ApiKey,
			"apiSecret": form.ApiSecret,
			"capital":   form.Capital,
			"status":    form.Status,
		}}
		if _, err = models.CustomerCollection.UpdateByID(
			context.TODO(),
			uId,
			update,
		); err != nil {
			resp.Error(err)
			return
		}

		resp.Success("Customer update successful", "/customer")
	})
}
