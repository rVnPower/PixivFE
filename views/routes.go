package views

import (
	"math"
	"net/http"
	"pixivfe/configs"
	"pixivfe/handler"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var PC *models.PixivClient

func artwork_page(c *gin.Context) {
	id := c.Param("id")
	illust, _ := PC.GetArtworkByID(id)
	related, _ := PC.GetRelatedArtworks(id)
	artist_info, _ := PC.GetUserInformation(illust.UserID, 1)

	c.HTML(http.StatusOK, "artwork.html", gin.H{
		"Illust":  illust,
		"Related": related,
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
	page, ok := c.GetQuery("page")

	if !ok {
		page = "1"
	}

	pageInt, _ := strconv.Atoi(page)
	user, _ := PC.GetUserInformation(id, pageInt)

	worksCount, _ := PC.GetUserArtworksCount(id)
	pageLimit := math.Ceil(float64(worksCount)/30) + 1.0

	c.HTML(http.StatusOK, "user.html", gin.H{"User": user, "PageLimit": int(pageLimit), "Page": pageInt})
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
}
