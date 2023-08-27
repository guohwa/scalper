package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"scalper/customers"
	forms "scalper/forms/customer"
	"scalper/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var customerHandler = &customerhandler{}

type customerhandler struct {
	base
}

func (handler *customerhandler) Handle(router *gin.Engine) {
	router.GET("/customer", func(ctx *gin.Context) {
		user, ok := getUser(ctx)
		if !ok {
			ctx.Redirect(http.StatusFound, "/account/login")
			return
		}

		filter := bson.M{
			"userId": user.ID,
		}

		count, err := models.CustomerCollection.CountDocuments(
			context.TODO(),
			filter,
			options.Count().SetMaxTime(2*time.Second))
		if err != nil {
			handler.Error(ctx, err)
			return
		}

		page, err := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
		if err != nil {
			handler.Error(ctx, err)
			return
		}
		limit, err := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)
		if err != nil {
			handler.Error(ctx, err)
			return
		}
		cursor, err := models.CustomerCollection.Find(
			context.TODO(),
			filter, options.Find().SetSkip((page-1)*limit).SetLimit(limit),
		)
		if err != nil {
			handler.Error(ctx, err)
			return
		}

		var items []models.Customer
		if err = cursor.All(context.TODO(), &items); err != nil {
			handler.Error(ctx, err)
			return
		}

		handler.HTML(ctx, "customer/index.html", Context{
			"page":  page,
			"count": count,
			"limit": limit,
			"items": items,
		})
	})

	router.GET("/customer/add", func(ctx *gin.Context) {
		handler.HTML(ctx, "customer/add.html", Context{})
	})

	router.POST("/customer/save", func(ctx *gin.Context) {
		user, ok := getUser(ctx)
		if !ok {
			ctx.Redirect(http.StatusFound, "/account/login")
			return
		}

		form := forms.Save{}

		if err := ctx.ShouldBind(&form); err != nil {
			handler.Error(ctx, err)
			return
		}

		customer := models.Customer{
			ID:        primitive.NewObjectID(),
			UserID:    user.ID,
			Name:      form.Name,
			ApiKey:    form.ApiKey,
			ApiSecret: form.ApiKey,
			Capital:   form.Capital,
			Position:  "",
			Status:    form.Status,
		}
		if _, err := models.CustomerCollection.InsertOne(
			context.TODO(),
			customer,
		); err != nil {
			handler.Error(ctx, err)
			return
		}

		customers.Add(customer)

		handler.Success(ctx, "Customer save successful", "/customer")
	})

	router.GET("/customer/edit/:id", func(ctx *gin.Context) {
		sId := ctx.Param("id")
		uId, err := primitive.ObjectIDFromHex(sId)
		if err != nil {
			handler.Error(ctx, err)
			return
		}

		customer := models.Customer{}
		if err := models.CustomerCollection.FindOne(context.TODO(), bson.M{
			"_id": uId,
		}).Decode(&customer); err != nil {
			handler.Error(ctx, err)
			return
		}

		handler.HTML(ctx, "customer/edit.html", Context{
			"item": customer,
		})
	})

	router.POST("/customer/update/:id", func(ctx *gin.Context) {
		form := forms.Update{}
		sId := ctx.Param("id")
		uId, err := primitive.ObjectIDFromHex(sId)
		if err != nil {
			handler.Error(ctx, err)
			return
		}

		if err := ctx.ShouldBind(&form); err != nil {
			handler.Error(ctx, err)
			return
		}

		var customer models.Customer
		update := bson.M{"$set": bson.M{
			"name":      form.Name,
			"apiKey":    form.ApiKey,
			"apiSecret": form.ApiSecret,
			"capital":   form.Capital,
			"status":    form.Status,
		}}
		if err = models.CustomerCollection.FindOneAndUpdate(
			context.TODO(),
			bson.M{
				"_id": uId,
			},
			update,
			options.FindOneAndUpdate().SetReturnDocument(options.After),
		).Decode(&customer); err != nil {
			handler.Error(ctx, err)
			return
		}

		customers.Set(customer)

		handler.Success(ctx, "Customer update successful", "/customer")
	})
}
