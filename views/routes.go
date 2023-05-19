package views

import (
	"net/http"
	"pixivfe/configs"
	"pixivfe/handler"
	"time"

	"github.com/gin-gonic/gin"
)

var PC *models.PixivClient

func artwork_page(c *gin.Context) {
	id := c.Param("id")
	illust, _ := PC.GetArtworkByID(id)
	related, _ := PC.GetRelatedArtworks(id)
	recent_by_artist, _ := PC.GetUserArtworks(illust.UserID)
	artist_info, _ := PC.GetUserInformation(illust.UserID)

	c.HTML(http.StatusOK, "artwork.html", gin.H{
		"Illust":  illust,
		"Related": related,
		"Recent":  recent_by_artist,
		"Artist":  artist_info,
	})
}

// func index_page(c *gin.Context) {
// 	recommended, _ := handler.GetRecommendedIllust(c)
// 	ranking, _ := handler.GetRankingIllust(c, "day")
// 	spotlight := handler.GetSpotlightArticle(c)
// 	newest, _ := handler.GetNewestIllust(c)
// 	c.HTML(http.StatusOK, "index.html", gin.H{
// 		"Recommended": recommended,
// 		"Rankings":    ranking,
// 		"Spotlights":  spotlight,
// 		"Newest":      newest,
// 	})
// }

func user_page(c *gin.Context) {
	id := c.Param("id")
	user, _ := PC.GetUserInformation(id)
	recent, _ := PC.GetUserArtworks(id)
	c.HTML(http.StatusOK, "user.html", gin.H{"User": user, "Recent": recent})
}

func getUserInformation(c *gin.Context) {
	id := c.Param("id")
	data, _ := PC.GetUserInformation(id)

	c.IndentedJSON(http.StatusOK, data)
}

func NewPixivClient(timeout int) *models.PixivClient {
	transport := &http.Transport{Proxy: http.ProxyFromEnvironment}
	client := &http.Client{
		Timeout:   time.Duration(timeout) * time.Millisecond,
		Transport: transport,
	}

	pc := &models.PixivClient{
		Client: client,
		Header: make(map[string]string),
		Cookie: make(map[string]string),
		Lang:   "en",
	}

	return pc
}

func SetupRoutes(r *gin.Engine) {
	PC = NewPixivClient(5000)
	PC.SetSessionID(configs.Configs.PHPSESSID)
	PC.SetUserAgent(configs.Configs.UserAgent)
	// r.GET("/", index_page)
	r.GET("artworks/:id", artwork_page)
	r.GET("users/:id", user_page)
	r.GET("api/users/:id", getUserInformation)
}
