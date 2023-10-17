package views

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"codeberg.org/vnpower/pixivfe/configs"
	"codeberg.org/vnpower/pixivfe/models"
	"github.com/gofiber/fiber/v2"
)

func get_session_value(c *fiber.Ctx, key string) *string {
	sess, err := configs.Store.Get(c)
	if err != nil {
		panic(err)
	}
	value := sess.Get(key)
	if value != nil {
		placeholder := value.(string)
		return &placeholder
	}
	return nil
}

func artwork_page(c *fiber.Ctx) error {
	image_proxy := get_session_value(c, "image-proxy")
	if image_proxy == nil {
		image_proxy = &configs.ProxyServer
	}

	id := c.Params("id")
	if _, err := strconv.Atoi(id); err != nil {
		return errors.New("Bad id")
	}

	illust, err := PC.GetArtworkByID(id)
	if err != nil {
		return err
	}

	illust.ProxyImages(*image_proxy)

	// Optimize this
	return c.Render("pages/artwork", fiber.Map{
		"Illust":   illust,
		"Title":    illust.Title,
		"PageType": "artwork",
	})
}

func index_page(c *fiber.Ctx) error {
	had_token := true
	image_proxy := get_session_value(c, "image-proxy")
	if image_proxy == nil {
		image_proxy = &configs.ProxyServer
	}
	token := get_session_value(c, "token")
	if token == nil {
		had_token = false
		token = &configs.Token
	}

	PC := NewPixivClient(5000)
	PC.SetSessionID(*token)
	PC.SetUserAgent(configs.UserAgent)
	PC.AddHeader("Accept-Language", "en-US,en;q=0.5")

	mode := c.Query("mode", "all")

	artworks, err := PC.GetLandingPage(mode)
	if err != nil {
		return err
	}

	if had_token {
		artworks.Following = models.ProxyShortArtworkSlice(artworks.Following, *image_proxy)
		artworks.Commissions = models.ProxyShortArtworkSlice(artworks.Commissions, *image_proxy)
		artworks.Recommended = models.ProxyShortArtworkSlice(artworks.Recommended, *image_proxy)
		artworks.Newest = models.ProxyShortArtworkSlice(artworks.Newest, *image_proxy)
		artworks.Users = models.ProxyShortArtworkSlice(artworks.Users, *image_proxy)
		artworks.RecommendByTags = models.ProxyRecommendedByTagsSlice(artworks.RecommendByTags, *image_proxy)
	}

	artworks.Rankings = models.ProxyShortArtworkSlice(artworks.Rankings, *image_proxy)
	artworks.Pixivision = models.ProxyPixivisionSlice(artworks.Pixivision, *image_proxy)

	return c.Render("pages/index", fiber.Map{"Title": "Landing", "Artworks": artworks, "Token": had_token})
}

func user_page(c *fiber.Ctx) error {
	image_proxy := get_session_value(c, "image-proxy")
	if image_proxy == nil {
		image_proxy = &configs.ProxyServer
	}

	id := c.Params("id")
	if _, err := strconv.Atoi(id); err != nil {
		return err
	}
	category := c.Params("category", "artworks")
	if !(category == "artworks" || category == "illustrations" || category == "manga" || category == "bookmarks") {
		return errors.New("Invalid work category: only illustrations, manga, artworks and bookmarks are available")
	}

	page := c.Query("page", "1")
	pageInt, _ := strconv.Atoi(page)

	user, err := PC.GetUserInformation(id, category, pageInt)
	if err != nil {
		return err
	}

	user.ProxyImages(*image_proxy)

	var worksCount int

	worksCount = user.ArtworksCount
	pageLimit := math.Ceil(float64(worksCount) / 30.0)

	return c.Render("pages/user", fiber.Map{"Title": user.Name, "User": user, "Category": category, "PageLimit": int(pageLimit), "Page": pageInt})
}

func ranking_page(c *fiber.Ctx) error {
	image_proxy := get_session_value(c, "image-proxy")
	if image_proxy == nil {
		image_proxy = &configs.ProxyServer
	}

	queries := make(map[string]string, 4)
	queries["Mode"] = c.Query("mode", "daily")
	queries["Content"] = c.Query("content", "all")
	queries["Date"] = c.Query("date", "")

	page := c.Query("page", "1")
	pageInt, _ := strconv.Atoi(page)

	response, err := PC.GetRanking(queries["Mode"], queries["Content"], queries["Date"], page)
	if err != nil {
		return err
	}

	response.ProxyImages(*image_proxy)

	return c.Render("pages/rank", fiber.Map{
		"Title":   "Ranking",
		"Data":    response,
		"Queries": queries,
		"Page":    pageInt,
	})
}

func newest_artworks_page(c *fiber.Ctx) error {
	image_proxy := get_session_value(c, "image-proxy")
	if image_proxy == nil {
		image_proxy = &configs.ProxyServer
	}

	worktype := c.Query("type", "illust")

	r18 := c.Query("r18", "false")

	works, err := PC.GetNewestArtworks(worktype, r18)
	if err != nil {
		return err
	}

	works = models.ProxyShortArtworkSlice(works, *image_proxy)

	return c.Render("pages/newest", fiber.Map{
		"Items": works,
		"Title": "Newest works",
	})
}

func search_page(c *fiber.Ctx) error {
	image_proxy := get_session_value(c, "image-proxy")
	if image_proxy == nil {
		image_proxy = &configs.ProxyServer
	}

	queries := make(map[string]string, 3)
	queries["Mode"] = c.Query("mode", "safe")
	queries["Category"] = c.Query("category", "artworks")
	queries["Order"] = c.Query("order", "date_d")

	name := c.Params("name")

	page := c.Query("page", "1")
	pageInt, _ := strconv.Atoi(page)

	tag, err := PC.GetTagData(name)
	if err != nil {
		return err
	}
	if len(tag.Metadata.Image) > 0 {
		tag.Metadata.Image = models.ProxyImage(tag.Metadata.Image, *image_proxy)
	}
	result, err := PC.GetSearch(queries["Category"], name, queries["Order"], queries["Mode"], page)
	if err != nil {
		return err
	}

	result.ProxyImages(*image_proxy)

	return c.Render("pages/tag", fiber.Map{"Title": "Results for " + tag.Name, "Tag": tag, "Data": result, "Queries": queries, "Page": pageInt})
}

func search(c *fiber.Ctx) error {
	name := c.FormValue("name")

	return c.Redirect("/tags/"+name, http.StatusFound)
}

func discovery_page(c *fiber.Ctx) error {
	image_proxy := get_session_value(c, "image-proxy")
	if image_proxy == nil {
		image_proxy = &configs.ProxyServer
	}

	mode := c.Query("mode", "safe")

	artworks, err := PC.GetDiscoveryArtwork(mode, 100)
	if err != nil {
		return err

	}
	artworks = models.ProxyShortArtworkSlice(artworks, *image_proxy)

	return c.Render("pages/discovery", fiber.Map{"Title": "Discovery", "Artworks": artworks})
}

func ranking_log_page(c *fiber.Ctx) error {
	image_proxy := get_session_value(c, "image-proxy")
	if image_proxy == nil {
		image_proxy = &configs.ProxyServer
	}

	mode := c.Query("mode", "daily")
	date := c.Query("date", "")

	var year int
	var month int
	var monthLit string

	// If the user supplied a date
	if len(date) == 6 {
		var err error
		year, err = strconv.Atoi(date[:4])
		if err != nil {
			return err
		}
		month, err = strconv.Atoi(date[4:])
		if err != nil {
			return err
		}
	} else {
		now := time.Now()
		year = now.Year()
		month = int(now.Month())
	}

	realDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	monthLit = realDate.Month().String()

	monthBefore := realDate.AddDate(0, -1, 0)
	monthAfter := realDate.AddDate(0, 1, 0)
	thisMonthLink := fmt.Sprintf("%d%02d", realDate.Year(), realDate.Month())
	monthBeforeLink := fmt.Sprintf("%d%02d", monthBefore.Year(), monthBefore.Month())
	monthAfterLink := fmt.Sprintf("%d%02d", monthAfter.Year(), monthAfter.Month())

	render, err := PC.GetRankingLog(mode, year, month, *image_proxy)
	if err != nil {
		return err
	}

	return c.Render("pages/ranking_log", fiber.Map{"Title": "Ranking calendar", "Render": render, "Mode": mode, "Month": monthLit, "Year": year, "MonthBefore": monthBeforeLink, "MonthAfter": monthAfterLink, "ThisMonth": thisMonthLink})
}

func following_works_page(c *fiber.Ctx) error {
	image_proxy := get_session_value(c, "image-proxy")
	if image_proxy == nil {
		image_proxy = &configs.ProxyServer
	}
	token := get_session_value(c, "token")
	if token == nil {
		return c.Redirect("/login")
	}
	queries := make(map[string]string, 2)
	queries["Mode"] = c.Query("mode", "all")
	queries["Page"] = c.Query("page", "1")
	pageInt, _ := strconv.Atoi(queries["Page"])

	artworks, err := PC.GetNewestFromFollowing(queries["Mode"], queries["Page"], *token)
	if err != nil {
		return err
	}
	artworks = models.ProxyShortArtworkSlice(artworks, *image_proxy)

	return c.Render("pages/following", fiber.Map{"Title": "Following works", "Queries": queries, "Artworks": artworks, "Page": pageInt})
}

func your_bookmark_page(c *fiber.Ctx) error {
	token := get_session_value(c, "token")
	if token == nil {
		return c.Redirect("/login")
	}

	// The left part of the token is the member ID
	userId := strings.Split(*token, "_")

	c.Redirect("/users/" + userId[0] + "/bookmarks#checkpoint")
	return nil
}

func login_page(c *fiber.Ctx) error {
	return c.Render("pages/login", fiber.Map{})
}

func settings_page(c *fiber.Ctx) error {
	return c.Render("pages/settings", fiber.Map{})
}

func settings_post(c *fiber.Ctx) error {
	t := c.Params("type")
	error := ""

	switch t {
	case "image_server":
		error = set_image_server(c)
	case "token":
		error = set_token(c)
	case "logout":
		error = set_logout(c)
	default:
		error = "No method available"
	}

	if error != "" {
		return errors.New(error)
	}
	c.Redirect("/settings")
	return nil
}

func get_logged_in_user(c *fiber.Ctx) error {
	token := get_session_value(c, "token")
	if token == nil {
		return c.Redirect("/login")
	}

	// The left part of the token is the member ID
	userId := strings.Split(*token, "_")

	c.Redirect("/users/" + userId[0])
	return nil
}
