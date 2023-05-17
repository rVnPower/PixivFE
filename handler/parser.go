package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pixivfe/entity"
	"regexp"

	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

func ImageProxy(url string) string {

	regex := regexp.MustCompile(`i\.pximg\.net`)
	proxy := "px2.rainchan.win"

	return regex.ReplaceAllString(url, proxy)
}

func ParseImages(data string) ([]entity.Image, error) {
	// Parse illusts images
	var images []entity.Image

	if gjson.Get(data, "meta_single_page.original_image_url").Exists() {
		var image entity.Image
		image.Small = ImageProxy(GetInnerJSON(data, "image_urls.square_medium"))
		image.Medium = ImageProxy(GetInnerJSON(data, "image_urls.medium"))
		image.Large = ImageProxy(GetInnerJSON(data, "image_urls.large"))
		image.Original = ImageProxy(GetInnerJSON(data, "meta_single_page.original_image_url"))

		images = append(images, image)
	} else {
		g := GetInnerJSON(data, "meta_pages.#.image_urls")

		err := json.Unmarshal([]byte(g), &images)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse JSON for images.\n%s", err)
		}
	}

	return images, nil
}

func ParseIllust(data string) (entity.Illust, error) {
	var illust entity.Illust

	images, err := ParseImages(data)

	if err != nil {
		return illust, fmt.Errorf("Failed to parse images from illust.\n%s", err)
	}

	illust.Images = images

	err = json.Unmarshal([]byte(data), &illust)
	if err != nil {
		return illust, fmt.Errorf("Failed to parse JSON for illust.\n%s", err)
	}

	return illust, nil
}

func ParseIllusts(data string) ([]entity.Illust, error) {
	var illusts []entity.Illust
	g := gjson.Get(data, "illusts").Array()

	for _, illust := range g {
		illust, err := ParseIllust(illust.String())

		if err != nil {
			return nil, fmt.Errorf("Failed to parse illusts.\n%s", err)
		}

		illusts = append(illusts, illust)
	}

	return illusts, nil
}

func GetRecommendedIllust(c *gin.Context) ([]entity.Illust, error) {
	URL := "https://hibi.cocomi.cf/api/pixiv/illust_recommended"

	s := Request(URL)
	illusts, err := ParseIllusts(s)

	if err != nil {
		return illusts, fmt.Errorf("Failed to get recommended illusts.\n%s", err)
	}

	return illusts, nil
}

func GetRankingIllust(c *gin.Context, mode string) ([]entity.Illust, error) {
	URL := "https://hibi.cocomi.cf/api/pixiv/rank?page=1&mode=" + mode

	s := Request(URL)
	illusts, err := ParseIllusts(s)

	if err != nil {
		return illusts, fmt.Errorf("Failed to get recommended illusts.\n%s", err)
	}

	return illusts, nil
}

func GetNewestIllust(c *gin.Context) ([]entity.Illust, error) {
	URL := "https://hibi.cocomi.cf/api/pixiv/illust_new"

	s := Request(URL)
	illusts, err := ParseIllusts(s)

	if err != nil {
		return illusts, fmt.Errorf("Failed to get recommended illusts.\n%s", err)
	}

	return illusts, nil
}

func GetMemberIllust(c *gin.Context, id string) ([]entity.Illust, error) {
	URL := "https://hibi.cocomi.cf/api/pixiv/member_illust?id=" + id

	s := Request(URL)
	illusts, err := ParseIllusts(s)

	if err != nil {
		return illusts, fmt.Errorf("Failed to get a member's illusts.\n%s", err)
	}

	return illusts, nil
}

func GetRelatedIllust(c *gin.Context) ([]entity.Illust, error) {
	id := c.Param("id")
	URL := "https://hibi.cocomi.cf/api/pixiv/related?id=" + id

	s := Request(URL)
	illusts, err := ParseIllusts(s)

	if err != nil {
		return illusts, fmt.Errorf("Failed to get related illusts.\n%s", err)
	}

	return illusts, nil
}

func GetIllustByID(c *gin.Context) (entity.Illust, error) {
	id := c.Param("id")
	URL := "https://hibi.cocomi.cf/api/pixiv/illust?id=" + id

	s := Request(URL)
	g := GetInnerJSON(s, "illust")
	illust, err := ParseIllust(g)

	if err != nil {
		return illust, fmt.Errorf("Failed to get illust by ID.\n%s", err)
	}

	return illust, nil
}

func GetSpotlightArticle(c *gin.Context) []entity.Spotlight {
	URL := "https://hibi.cocomi.cf/api/pixiv/spotlights?lang=en"
	// URL := "https://now.pixiv.pics/api/pixivision?lang=en"
	var articles []entity.Spotlight

	s := Request(URL)
	g := GetInnerJSON(s, "spotlight_articles")

	err := json.Unmarshal([]byte(g), &articles)
	if err != nil {
		panic("Failed to parse JSON")
	}

	return articles
}

func GetUserInfo(c *gin.Context) (entity.User, error) {
	id := c.Param("id")
	URL := "https://hibi.cocomi.cf/api/pixiv/member?id=" + id
	var user entity.User

	s := Request(URL)
	user_string := GetInnerJSON(s, "user")
	profile_string := GetInnerJSON(s, "profile")

	err := json.Unmarshal([]byte(user_string), &user)
	if err != nil {
		panic("Failed to parse JSON")
	}

	err = json.Unmarshal([]byte(profile_string), &user)
	if err != nil {
		panic("Failed to parse JSON")
	}

	return user, nil
}

func GetInnerJSON(resp string, key string) string {
	// As I see, the API always start its response with { "key": [...] }
	return gjson.Get(resp, key).String()
}

func Request(URL string) string {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", URL, nil)
	req.Header.Set("accept-language", "en")
	resp, err := client.Do(req)

	if err != nil {
		panic("Failed to make a request to " + URL)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic("Failed to parse response")
	}

	return string(body)
}
