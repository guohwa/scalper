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

var paramHandler = &paramhandler{}

type paramhandler struct {
}

func (handler *paramhandler) Handle(router *gin.Engine) {
	router.GET("/param", func(ctx *gin.Context) {
		resp := response.New(ctx)
		resp.HTML("param/index.html", response.Context{
			"param": config.Param,
		})
	})

	router.POST("/param", func(ctx *gin.Context) {
		resp := response.New(ctx)
		form := forms.Param{}

		if err := ctx.ShouldBind(&form); err != nil {
			resp.Error(err)
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
			"superTrend": bson.M{
				"demaLength": form.SuperTrendDemaLength,
				"atrLength":  form.SuperTrendAtrLength,
				"atrMult":    form.SuperTrendAtrMult,
			},
			"tutci": bson.M{
				"entry": form.TutciEntry,
			},
			"ssl": bson.M{
				"enable": form.SSLEnable,
				"length": form.SSLLength,
			},
			"pv": bson.M{
				"enable":    form.PVEnable,
				"threshold": form.PVThreshold,
			},
			"tsl": bson.M{
				"trailProfit":   form.TSLTrailProfit,
				"trailOffset":   form.TSLTrailOffset,
				"trailStopLoss": form.TSLStopLoss,
			},
		}}
		if _, err := models.ParamCollection.UpdateByID(
			context.TODO(),
			config.Param.ID,
			update,
		); err != nil {
			resp.Error(err)
			return
		}

		config.Param.Symbol.Precision.Price = form.SymbolPricePrecision
		config.Param.Symbol.Precision.Quantity = form.SymbolQuantityPrecision
		config.Param.SuperTrend.DemaLength = form.SuperTrendDemaLength
		config.Param.SuperTrend.AtrLength = form.SuperTrendAtrLength
		config.Param.SuperTrend.AtrMult = form.SuperTrendAtrMult
		config.Param.TuTCI.Entry = form.TutciEntry
		config.Param.SSL.Enable = form.SSLEnable
		config.Param.SSL.Length = form.SSLLength
		config.Param.PV.Enable = form.PVEnable
		config.Param.PV.Threshold = form.PVThreshold
		config.Param.TSL.TrailProfit = form.TSLTrailProfit
		config.Param.TSL.TrailOffset = form.TSLTrailOffset
		config.Param.TSL.StopLoss = form.TSLStopLoss

		resp.Success("Param update successful", "")
	})
}
