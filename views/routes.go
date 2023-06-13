package views

import (
	"errors"
	"math"
	"net/http"
	"pixivfe/configs"
	"pixivfe/handler"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

var PC *handler.PixivClient

func artwork_page(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := strconv.Atoi(id); err != nil {
		return errors.New("Bad id")
	}

	illust, err := PC.GetArtworkByID(id)
	if err != nil {
		return err
	}

	related, _ := PC.GetRelatedArtworks(id)
	comments, _ := PC.GetArtworkComments(id)
	artist_info, err := PC.GetUserInformation(illust.UserID, 1)

	if err != nil {
		return err
	}

	return c.Render("artwork", fiber.Map{
		"Illust":   illust,
		"Related":  related,
		"Artist":   artist_info,
		"Comments": comments,
		"Title":    illust.Title,
	}, "layout")
}

func index_page(c *fiber.Ctx) error {
	// recommended, _ := handler.GetRecommendedIllust(c)
	// ranking, _ := handler.GetRankingIllust(c, "day")
	// spotlight := handler.GetSpotlightArticle(c)
	// newest, _ := handler.GetNewestIllust(c)
	// return c.Render(http.StatusOK, "index.html", fiber.Map{
	// 	"Recommended": recommended,
	// 	"Rankings":    ranking,
	// 	"Spotlights":  spotlight,
	// 	"Newest":      newest,
	// })
	return c.Render("temp", fiber.Map{"Title": "Test"}, "layout")
}

func user_page(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := strconv.Atoi(id); err != nil {
		return err
	}
	page := c.Query("page", "1")

	pageInt, _ := strconv.Atoi(page)
	user, err := PC.GetUserInformation(id, pageInt)
	if err != nil {
		return err
	}

	worksCount, _ := PC.GetUserArtworksCount(id)
	pageLimit := math.Ceil(float64(worksCount)/30.0) + 1.0

	return c.Render("user", fiber.Map{"Title": user.Name, "User": user, "PageLimit": int(pageLimit), "Page": pageInt}, "layout")
}

func ranking_page(c *fiber.Ctx) error {
	mode := c.Query("mode", "daily")

	content := c.Query("content", "all")

	page := c.Query("page", "1")

	pageInt, _ := strconv.Atoi(page)

	response, err := PC.GetRanking(mode, content, page)
	if err != nil {
		return err
	}

	return c.Render("rank", fiber.Map{
		"Title":   "Ranking",
		"Items":   response.Artworks,
		"Mode":    mode,
		"Content": content,
		"Page":    pageInt}, "layout")
}

func newest_artworks_page(c *fiber.Ctx) error {
	worktype := c.Query("type", "illust")

	r18 := c.Query("r18", "false")

	works, err := PC.GetNewestArtworks(worktype, r18)
	if err != nil {
		return err
	}

	return c.Render("newest", fiber.Map{
		"Items": works,
		"Title": "Newest works",
	}, "layout")
}

func search_page(c *fiber.Ctx) error {
	name := c.Params("name")

	page := c.Query("page", "1")

	order := c.Query("order", "date_d")

	mode := c.Query("mode", "safe")

	category := c.Query("category", "artworks")

	tag, err := PC.GetTagData(name)
	if err != nil {
		return err
	}
	result, err := PC.GetSearch(category, name, order, mode, page)
	if err != nil {
		return err
	}

	queries := map[string]string{
		"Page":     page,
		"Order":    order,
		"Mode":     mode,
		"Category": category,
	}
	return c.Render("tag", fiber.Map{"Title": "Results for " + tag.Name, "Tag": tag, "Data": result, "Queries": queries}, "layout")
}

func search(c *fiber.Ctx) error {
	name := c.FormValue("name")

	return c.Redirect("/tags/"+name, http.StatusFound)
}

func discovery_page(c *fiber.Ctx) error {
	mode := c.Query("mode", "safe")

	println("Start")
	artworks, err := PC.GetDiscoveryArtwork(mode, 100)
	println("End")
	if err != nil {
		return err
	}

	return c.Render("discovery", fiber.Map{"Title": "Discovery", "Artworks": artworks}, "layout")
}

// func not_found_page(c *fiber.Ctx) {
// 	return c.Render(http.StatusNotFound, "error.html", fiber.Map{
// 		"Title": "Not found",
// 		"Error": "Route " + c.Request.URL.Path + " not found.",
// 	})
// }

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

func SetupRoutes(r *fiber.App) {
	PC = NewPixivClient(5000)
	PC.SetSessionID(configs.Token)
	PC.SetUserAgent(configs.UserAgent)

	r.Get("/", index_page)
	r.Get("artworks/:id", artwork_page)
	r.Get("users/:id", user_page)
	r.Get("newest", newest_artworks_page)
	r.Get("ranking", ranking_page)
	r.Get("tags/:name", search_page)
	r.Get("discovery", discovery_page)
	r.Post("tags", search)

	// 404 page
	// r.NoRoute(not_found_page)
}
