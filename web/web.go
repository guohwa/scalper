package web

import (
	"log"

	"scalper/config"
	"scalper/models"
	"scalper/web/handlers"
	"scalper/web/middleware"
	"scalper/web/middleware/pongo2gin"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/mongo/mongodriver"
	"github.com/gin-gonic/gin"
)

func Start() {
	if config.App.Mode == "PROD" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.SetTrustedProxies(config.App.Trust)

	router.Static("/static", "./web/static")
	router.StaticFile("/favicon.ico", "./web/static/favicon.ico")
	router.StaticFile("/robots.txt", "./web/static/robots.txt")

	store := mongodriver.NewStore(models.SessionCollection, 3600, true, []byte("!PassW0rd!"))
	router.Use(sessions.Sessions("session", store))

	router.HTMLRender = pongo2gin.New(pongo2gin.RenderOptions{
		TemplateDir: "./web/templates",
		ContentType: "text/html",
	})

	router.Use(middleware.Global())
	router.Use(middleware.Intercept(auth))

	handlers.Handle(router)

	log.Println("Server listen on " + config.App.Listen)
	router.Run(config.App.Listen)
}
