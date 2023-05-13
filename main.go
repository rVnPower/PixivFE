package main

import (
	"html/template"
	"net/http"
	"pixivfe/handler"

	"github.com/gin-gonic/gin"
)

func inc(n int) int {
	return n + 1
}

func index_page(c *gin.Context) {
	recommended := handler.GetRecommendedIllust(c)
	ranking := handler.GetRankingIllust(c, "")
	spotlight := handler.GetSpotlightArticle(c)
	c.HTML(http.StatusOK, "index.html", gin.H{"Recommended": recommended, "Rankings": ranking, "Spotlights": spotlight})
}

func main() {
	server := gin.Default()
	server.SetFuncMap(template.FuncMap{
		"inc": inc,
	})
	server.Static("css/", "./template/css")
	server.LoadHTMLGlob("template/*.html")
	server.GET("/", index_page)

	server.Run(":8000")
}
