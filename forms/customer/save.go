package forms

type Save struct {
	Name      string `form:"name" binding:"required"`
	ApiKey    string `form:"apiKey" binding:"required"`
	ApiSecret string `form:"apiSecret" binding:"required"`
	Status    string `form:"status" binding:"required,oneof=Enable Disable"`
}