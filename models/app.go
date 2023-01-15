package models

type App struct {
	Title  string   `json:"title"`
	Mode   string   `json:"mode"`
	Listen string   `json:"listen"`
	Trust  []string `json:"trust"`
}

func (app *App) Default() {
	app.Title = "Scalper"
	app.Mode = "Dev"
	app.Listen = ":8080"
	app.Trust = []string{
		"127.0.0.1",
	}
}
