package forms

type Param struct {
	SymbolPricePrecision    int     `form:"symbolPricePrecision" binding:"required"`
	SymbolQuantityPrecision int     `form:"symbolQuantityPrecision" binding:"required"`
	SuperTrendDemaLength    int     `form:"superTrendDemaLength" binding:"required"`
	SuperTrendAtrLength     int     `form:"superTrendAtrLength" binding:"required"`
	SuperTrendAtrMult       float64 `form:"superTrendAtrMult" binding:"required"`
	TutciEntry              int     `form:"tutciEntry" binding:"required"`
	SSLLength               int     `form:"sslLength" binding:"required"`
	PVThreshold             float64 `form:"pvThreshold" binding:"required"`
	TSLTrailProfit          float64 `form:"tslTrailProfit" binding:"required"`
	TSLTrailOffset          float64 `form:"tslTrailOffset" binding:"required"`
	TSLStopLoss             float64 `form:"tslStopLoss" binding:"required"`
}
