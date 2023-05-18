package main

import (
	"html/template"
	"pixivfe/configs"
	"pixivfe/views"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	server := gin.Default()

	// For rankings to increment a number by 1
	server.SetFuncMap(template.FuncMap{
		"inc": func(n int) int {
			return n + 1
		},
	})

	// Static files
	server.StaticFile("/favicon.ico", "./template/favicon.ico")
	server.Static("css/", "./template/css")

	// HTML templates, automatically loaded
	server.LoadHTMLGlob("template/*.html")

	// Routes/Views
	views.SetupRoutes(server)

	return server
}

func main() {
	configs.Configs.ReadConfig()
	println(configs.Configs.PHPSESSID)

	r := setupRouter()

	r.Run(":8080")
}
