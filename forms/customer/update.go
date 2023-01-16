package forms

type Update struct {
	Name      string  `form:"name" binding:"required"`
	ApiKey    string  `form:"apiKey" binding:"required"`
	ApiSecret string  `form:"apiSecret" binding:"required"`
	Capital   float64 `form:"capital" binding:"required"`
	Status    string  `form:"status" binding:"required,oneof=Enable Disable"`
}
