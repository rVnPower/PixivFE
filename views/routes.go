package views

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pixivfe/handler"
	"strconv"
)

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

func SetupRoutes(r *gin.Engine) {
	r.GET("/", index_page)
	r.GET("artworks/:id", artwork_page)
}
