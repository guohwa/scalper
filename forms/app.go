package forms

type App struct {
	Title string `form:"title" binding:"required"`
	Mode  string `form:"mode" binding:"required,oneof=debug test release"`
	Level string `form:"level" binding:"required,oneof=trace debug info warn error fatal panic"`
}
