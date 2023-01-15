package models

type Config struct {
	Symbol struct {
		Name      string `bson:"name"`
		Period    string `bson:"period"`
		Limit     int    `bson:"limit"`
		Precision struct {
			Price    int `bson:"price"`
			Quantity int `bson:"quantity"`
		} `bson:"precision"`
	}
	SuperTrend struct {
		DemaLength int     `bson:"demaLength"`
		AtrLength  int     `bson:"atrLength"`
		AtrMult    float64 `bson:"atrMult"`
	} `bson:"superTrend"`
	TuTCI struct {
		Entry int `bson:"entry"`
	} `bson:"tutci"`
	SSL struct {
		Length int `bson:"length"`
	} `bson:"ssl"`
	PV struct {
		Threshold float64 `bson:"threshold"`
	} `bson:"pv"`
	TSL struct {
		TrailProfit float64 `bson:"trailProfit"`
		TrailOffset float64 `bson:"trailOffset"`
		StopLoss    float64 `bson:"stopLoss"`
	} `bson:"tsl"`
}
