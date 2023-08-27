package handlers

import (
	"context"

	"scalper/config"
	"scalper/forms"
	"scalper/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

var paramHandler = &paramhandler{}

type paramhandler struct {
	base
}

func (handler *paramhandler) Handle(router *gin.Engine) {
	router.GET("/param", func(ctx *gin.Context) {
		handler.HTML(ctx, "param/index.html", Context{
			"param": config.Param,
		})
	})

	router.POST("/param", func(ctx *gin.Context) {
		form := forms.Param{}

		if err := ctx.ShouldBind(&form); err != nil {
			handler.Error(ctx, err)
			return
		}

		update := bson.M{"$set": bson.M{
			"symbol": bson.M{
				"name":   config.Param.Symbol.Name,
				"period": config.Param.Symbol.Period,
				"limit":  config.Param.Symbol.Limit,
				"precision": bson.M{
					"price":    form.SymbolPricePrecision,
					"quantity": form.SymbolQuantityPrecision,
				},
			},
			"ssl": bson.M{
				"length": form.SSLLength,
			},
			"tsl": bson.M{
				"trailProfit": form.TSLTrailProfit,
				"trailOffset": form.TSLTrailOffset,
				"stopLoss":    form.TSLStopLoss,
			},
		}}
		if _, err := models.ParamCollection.UpdateByID(
			context.TODO(),
			config.Param.ID,
			update,
		); err != nil {
			handler.Error(ctx, err)
			return
		}

		config.Param.Symbol.Precision.Price = form.SymbolPricePrecision
		config.Param.Symbol.Precision.Quantity = form.SymbolQuantityPrecision
		config.Param.SSL.Length = form.SSLLength
		config.Param.TSL.TrailProfit = form.TSLTrailProfit
		config.Param.TSL.TrailOffset = form.TSLTrailOffset
		config.Param.TSL.StopLoss = form.TSLStopLoss

		handler.Success(ctx, "Param update successful", "")
	})
}
