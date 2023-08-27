package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"scalper/log"
	"scalper/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var orderHandler = &orderhandler{}

type orderhandler struct {
	base
}

func (handler *orderhandler) Handle(router *gin.Engine) {
	router.GET("/order/*id", func(ctx *gin.Context) {
		user, ok := getUser(ctx)
		if !ok {
			ctx.Redirect(http.StatusFound, "/account/login")
			return
		}

		filter := bson.M{
			"userId": user.ID,
			"status": "Enable",
		}
		cursor, err := models.CustomerCollection.Find(
			context.TODO(),
			filter, options.Find(),
		)
		if err != nil {
			handler.Error(ctx, err)
			return
		}
		var customers []models.Customer
		if err = cursor.All(context.TODO(), &customers); err != nil {
			handler.Error(ctx, err)
			return
		}

		sId := strings.TrimLeft(ctx.Param("id"), "/")
		session := sessions.Default(ctx)
		if sId == "" {
			cId := session.Get("customer-id")
			if cId != nil {
				if id, ok := cId.(string); ok {
					sId = id
				}
			}
		} else {
			session.Set("customer-id", sId)
			if err := session.Save(); err != nil {
				log.Error(err)
			}
		}

		var customer models.Customer
		if sId != "" {
			for _, item := range customers {
				if item.ID.Hex() == sId {
					customer = item
				}
			}
		} else {
			if len(customers) > 0 {
				customer = customers[0]
			}
		}

		filter = bson.M{
			"customerId": customer.ID,
		}
		count, err := models.OrderCollection.CountDocuments(
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
		cursor, err = models.OrderCollection.Find(
			context.TODO(),
			filter,
			options.Find().SetSort(bson.M{"entryTime": -1}).SetSkip((page-1)*limit).SetLimit(limit),
		)
		if err != nil {
			handler.Error(ctx, err)
			return
		}

		var orders []models.Order
		if err = cursor.All(context.TODO(), &orders); err != nil {
			handler.Error(ctx, err)
			return
		}

		handler.HTML(ctx, "order/index.html", Context{
			"page":      page,
			"count":     count,
			"limit":     limit,
			"orders":    orders,
			"customer":  customer,
			"customers": customers,
		})
	})
}
