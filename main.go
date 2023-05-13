package main

import (
	"net/http"
	"pixivfe/handler"

	"github.com/gin-gonic/gin"
)

func test(c *gin.Context) {
	data := handler.GetRecommendedIllust(c)
	c.HTML(http.StatusOK, "index.html", data)
}

func main() {
	server := gin.Default()
	server.Static("css/", "./template/css")
	server.LoadHTMLGlob("template/*.html")
	// Listen and Servr in 0.0.0.0:8080
	server.GET("/", test)

	server.Run(":8080")
}
