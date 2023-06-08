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

var PC *handler.PixivClient

func artwork_page(c *gin.Context) {
	id := c.Param("id")
	if _, err := strconv.Atoi(id); err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Title": "Bad request",
			"Error": "Invalid artwork ID",
		})
		return
	}

	illust, err := PC.GetArtworkByID(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Title": "An error occured",
			"Error": err,
		})
		return
	}

	related, _ := PC.GetRelatedArtworks(id)
	comments, _ := PC.GetArtworkComments(id)
	artist_info, err := PC.GetUserInformation(illust.UserID, 1)

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Title": "An error occured",
			"Error": err,
		})
	}

	c.HTML(http.StatusOK, "artwork.html", gin.H{
		"Illust":   illust,
		"Related":  related,
		"Artist":   artist_info,
		"Comments": comments,
		"Title":    illust.Title,
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
	if _, err := strconv.Atoi(id); err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Title": "Bad request",
			"Error": "Invalid user ID",
		})
		return
	}
	page, ok := c.GetQuery("page")

	if !ok {
		page = "1"
	}

	pageInt, _ := strconv.Atoi(page)
	user, err := PC.GetUserInformation(id, pageInt)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Title": "An error occured",
			"Error": err,
		})
		return
	}

	worksCount, _ := PC.GetUserArtworksCount(id)
	pageLimit := math.Ceil(float64(worksCount)/30.0) + 1.0

	c.HTML(http.StatusOK, "user.html", gin.H{"Title": user.Name, "User": user, "PageLimit": int(pageLimit), "Page": pageInt})
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

	response, err := PC.GetRanking(mode, content, page)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Title": "An error occured",
			"Error": err,
		})
		return
	}

	c.HTML(http.StatusOK, "rank.html", gin.H{
		"Title":   "Ranking",
		"Items":   response.Artworks,
		"Mode":    mode,
		"Content": content,
		"Page":    pageInt})
}

func newest_artworks_page(c *gin.Context) {
	worktype, ok := c.GetQuery("type")

	if !ok {
		worktype = "illust"
	}
	r18, ok := c.GetQuery("r18")

	if !ok {
		r18 = "false"
	}

	works, err := PC.GetNewestArtworks(worktype, r18)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Title": "An error occured",
			"Error": err,
		})
		return
	}

	c.HTML(http.StatusOK, "newest.html", gin.H{
		"Items": works,
		"Title": "Newest works",
	})
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

	tag, err := PC.GetTagData(name)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Title": "An error occured",
			"Error": err,
		})
		return
	}
	result, _ := PC.GetSearch(category, name, order, mode, page)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Title": "An error occured",
			"Error": err,
		})
		return
	}

	queries := map[string]string{
		"Page":     page,
		"Order":    order,
		"Mode":     mode,
		"Category": category,
	}
	c.HTML(http.StatusOK, "tag.html", gin.H{"Title": "Results for " + tag.Name, "Tag": tag, "Data": result, "Queries": queries})
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

	artworks, err := PC.GetDiscoveryArtwork(mode, 300)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Title": "An error occured",
			"Error": err,
		})
		return
	}

	c.HTML(http.StatusOK, "discovery.html", gin.H{"Title": "Discovery", "Artworks": artworks})
}

func not_found_page(c *gin.Context) {
	c.HTML(http.StatusNotFound, "error.html", gin.H{
		"Title": "Not found",
		"Error": "Route " + c.Request.URL.Path + " not found.",
	})
}

func NewPixivClient(timeout int) *handler.PixivClient {
	transport := &http.Transport{Proxy: http.ProxyFromEnvironment}
	client := &http.Client{
		Timeout:   time.Duration(timeout) * time.Millisecond,
		Transport: transport,
	}

	pc := &handler.PixivClient{
		Client: client,
		Header: make(map[string]string),
		Cookie: make(map[string]string),
		Lang:   "en",
	}

	return pc
}

func SetupRoutes(r *gin.Engine) {
	PC = NewPixivClient(5000)
	PC.SetSessionID(configs.Token)
	PC.SetUserAgent(configs.UserAgent)
	r.GET("/", index_page)
	r.GET("artworks/:id", artwork_page)
	r.GET("users/:id", user_page)
	r.GET("newest", newest_artworks_page)
	r.GET("ranking", ranking_page)
	r.GET("tags/:name", search_page)
	r.GET("discovery", discovery_page)
	r.POST("tags", search)

	// 404 page
	r.NoRoute(not_found_page)
}
