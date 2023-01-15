package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Precision struct {
	Price    int `bson:"price"`
	Quantity int `bson:"quantity"`
}

type Symbol struct {
	Name      string    `bson:"name"`
	Period    string    `bson:"period"`
	Limit     int       `bson:"limit"`
	Precision Precision `bson:"precision"`
}

type SuperTrend struct {
	DemaLength int     `bson:"demaLength"`
	AtrLength  int     `bson:"atrLength"`
	AtrMult    float64 `bson:"atrMult"`
}

type TuTCI struct {
	Entry int `bson:"entry"`
}

type SSL struct {
	Length int `bson:"length"`
}

type PV struct {
	Threshold float64 `bson:"threshold"`
}

type TSL struct {
	TrailProfit float64 `bson:"trailProfit"`
	TrailOffset float64 `bson:"trailOffset"`
	StopLoss    float64 `bson:"stopLoss"`
}

type Param struct {
	ID         primitive.ObjectID `bson:"_id"`
	Symbol     Symbol             `bson:"symbol"`
	SuperTrend SuperTrend         `bson:"superTrend"`
	TuTCI      TuTCI              `bson:"tutci"`
	SSL        SSL                `bson:"ssl"`
	PV         PV                 `bson:"pv"`
	TSL        TSL                `bson:"tsl"`
}

func (param *Param) Default() {
	param.Symbol = Symbol{
		Name:   "ETHUSDT",
		Period: "15m",
		Limit:  1500,
		Precision: Precision{
			Price:    2,
			Quantity: 3,
		},
	}
	param.SuperTrend = SuperTrend{
		DemaLength: 9,
		AtrLength:  3,
		AtrMult:    5.0,
	}
	param.TuTCI = TuTCI{
		Entry: 4,
	}
	param.SSL = SSL{
		Length: 160,
	}
	param.PV = PV{
		Threshold: 2.0,
	}
	param.TSL = TSL{
		TrailProfit: 0.3,
		TrailOffset: 0.03,
		StopLoss:    2.5,
	}
}
