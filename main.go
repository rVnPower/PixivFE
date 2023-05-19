package main

import (
	"html/template"
	"pixivfe/configs"
	"pixivfe/views"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	server := gin.Default()

	server.SetFuncMap(template.FuncMap{
		"inc": func(n int) int {
			// For rankings to increment a number by 1
			return n + 1
		},
		"add": func(a int, b int) int {
			return a + b
		},

		"dec": func(n int) int {
			return n - 1
		},

		"isEmpty": func(s string) bool {
			return len(s) < 1
		},

		"isEmphasize": func(s string) bool {
			switch s {
			case
				"R-18",
				"R-18G":
				return true
			}
			return false
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

	r := setupRouter()

	r.Run(":8080")
}
