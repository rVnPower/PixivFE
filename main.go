package main

import (
	"html/template"
	"net/http"
	"pixivfe/handler"
	"strconv"

	"github.com/gin-gonic/gin"
)

func inc(n int) int {
	return n + 1
}

func artwork_page(c *gin.Context) {
	illust := handler.GetIllustByID(c)
	related := handler.GetRelatedIllust(c)
	recent_by_artist := handler.GetMemberIllust(c, strconv.Itoa(illust.Artist.ID))
	c.HTML(http.StatusOK, "artwork.html", gin.H{
		"Illust":  illust,
		"Related": related,
		"Recent":  recent_by_artist,
	})
}

func index_page(c *gin.Context) {
	recommended := handler.GetRecommendedIllust(c)
	ranking := handler.GetRankingIllust(c, "day")
	spotlight := handler.GetSpotlightArticle(c)
	newest := handler.GetNewestIllust(c)
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Recommended": recommended,
		"Rankings":    ranking,
		"Spotlights":  spotlight,
		"Newest":      newest,
	})
}

func main() {
	server := gin.Default()
	server.SetFuncMap(template.FuncMap{
		"inc": inc,
	})
	server.StaticFile("/favicon.ico", "./template/favicon.ico")
	server.Static("css/", "./template/css")
	server.LoadHTMLGlob("template/*.html")
	server.GET("/", index_page)
	server.GET("artworks/:id", artwork_page)

	server.Run(":8080")
}
