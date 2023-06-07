package main

import (
	"pixivfe/configs"
	"pixivfe/handler"
	"pixivfe/views"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	server := gin.Default()

	server.SetFuncMap(handler.GetTemplateFunctions())

	// Static files
	server.StaticFile("/favicon.ico", "./template/favicon.ico")
	server.Static("css/", "./template/css")
	server.Static("assets/", "./template/assets")

	// HTML templates, automatically loaded
	server.LoadHTMLGlob("template/*.html")

	// Routes/Views
	views.SetupRoutes(server)

	// Disable trusted proxies since we do not use any for now
	server.SetTrustedProxies(nil)

	return server
}

func main() {
	err := configs.ParseConfig()

	if err != nil {
		panic(err)
	}

	r := setupRouter()

	r.Run(":" + configs.Port)
}
