package views

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pixivfe/handler"
	"strconv"
)

func artwork_page(c *gin.Context) {
	illust, _ := handler.GetIllustByID(c)
	related, _ := handler.GetRelatedIllust(c)
	recent_by_artist, _ := handler.GetMemberIllust(c, strconv.Itoa(illust.Artist.ID))
	c.HTML(http.StatusOK, "artwork.html", gin.H{
		"Illust":  illust,
		"Related": related,
		"Recent":  recent_by_artist,
	})
}

func index_page(c *gin.Context) {
	recommended, _ := handler.GetRecommendedIllust(c)
	ranking, _ := handler.GetRankingIllust(c, "day")
	spotlight := handler.GetSpotlightArticle(c)
	newest, _ := handler.GetNewestIllust(c)
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Recommended": recommended,
		"Rankings":    ranking,
		"Spotlights":  spotlight,
		"Newest":      newest,
	})
}

func user_page(c *gin.Context) {
	user, _ := handler.GetUserInfo(c)
	recent, _ := handler.GetMemberIllust(c, c.Param("id"))
	c.HTML(http.StatusOK, "user.html", gin.H{"User": user, "Recent": recent})
}

func SetupRoutes(r *gin.Engine) {
	r.GET("/", index_page)
	r.GET("artworks/:id", artwork_page)
	r.GET("user/:id", user_page)
}
