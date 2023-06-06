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
	comments, _ := PC.GetArtworkComments(id)
	artist_info, _ := PC.GetUserInformation(illust.UserID, 1)

	c.HTML(http.StatusOK, "artwork.html", gin.H{
		"Illust":   illust,
		"Related":  related,
		"Artist":   artist_info,
		"Comments": comments,
	})
}

func index_page(c *gin.Context) {
	// recommended, _ := handler.GetRecommendedIllust(c)
	// ranking, _ := handler.GetRankingIllust(c, "day")
	// spotlight := handler.GetSpotlightArticle(c)
	// newest, _ := handler.GetNewestIllust(c)
	// c.HTML(http.StatusOK, "index.html", gin.H{
	// 	"Recommended": recommended,
	// 	"Rankings":    ranking,
	// 	"Spotlights":  spotlight,
	// 	"Newest":      newest,
	// })
	c.HTML(http.StatusOK, "temp.html", gin.H{})
}

func user_page(c *gin.Context) {
	id := c.Param("id")
	page, ok := c.GetQuery("page")

	if !ok {
		page = "1"
	}

	pageInt, _ := strconv.Atoi(page)
	user, _ := PC.GetUserInformation(id, pageInt)

	worksCount, _ := PC.GetUserArtworksCount(id)
	pageLimit := math.Ceil(float64(worksCount)/float64(configs.Configs.PageItems)) + 1.0

	c.HTML(http.StatusOK, "user.html", gin.H{"User": user, "PageLimit": int(pageLimit), "Page": pageInt})
}

func ranking_page(c *gin.Context) {
	mode, ok := c.GetQuery("mode")

	if !ok {
		mode = "daily"
	}

	content, ok := c.GetQuery("content")

	if !ok {
		content = "all"
	}

	page, ok := c.GetQuery("page")

	if !ok {
		page = "1"
	}

	pageInt, _ := strconv.Atoi(page)

	response, _ := PC.GetRanking(mode, content, page)

	c.HTML(http.StatusOK, "rank.html", gin.H{"Items": response.Artworks,
		"Mode":    mode,
		"Content": content,
		"Page":    pageInt})
}

func newestArtworksPage(c *gin.Context) {
	worktype, ok := c.GetQuery("type")

	if !ok {
		worktype = "illust"
	}
	r18, ok := c.GetQuery("r18")

	if !ok {
		r18 = "false"
	}

	works, _ := PC.GetNewestArtworks(worktype, r18)

	c.HTML(http.StatusOK, "newest.html", gin.H{"Items": works})
}

func search_page(c *gin.Context) {
	name := c.Param("name")

	page, ok := c.GetQuery("page")

	if !ok {
		page = "1"
	}

	order, ok := c.GetQuery("order")

	if !ok {
		order = "date_d"
	}

	mode, ok := c.GetQuery("mode")

	if !ok {
		mode = "safe"
	}

	category, ok := c.GetQuery("category")

	if !ok {
		category = "artworks"
	}

	tag, _ := PC.GetTagData(name)
	result, _ := PC.GetSearch(category, name, order, mode, page)

	queries := map[string]string{
		"Page":     page,
		"Order":    order,
		"Mode":     mode,
		"Category": category,
	}
	c.HTML(http.StatusOK, "tag.html", gin.H{"Tag": tag, "Data": result, "Queries": queries})
}

func search(c *gin.Context) {
	name := c.PostForm("name")

	c.Redirect(http.StatusFound, "/tags/"+name)
}

func discovery_page(c *gin.Context) {
	mode, ok := c.GetQuery("mode")

	if !ok {
		mode = "safe"
	}

	artworks, _ := PC.GetDiscoveryArtwork(mode)
	c.HTML(http.StatusOK, "discovery.html", gin.H{"Artworks": artworks})
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
	r.GET("/", index_page)
	r.GET("artworks/:id", artwork_page)
	r.GET("users/:id", user_page)
	r.GET("newest", newestArtworksPage)
	r.GET("ranking", ranking_page)
	r.GET("tags/:name", search_page)
	r.GET("discovery", discovery_page)
	r.POST("tags", search)
}
