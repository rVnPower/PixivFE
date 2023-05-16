package handler

import (
	"encoding/json"
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

func ParseImages(data string) []entity.Image {
	var images []entity.Image

	if gjson.Get(data, "meta_single_page.original_image_url").Exists() {
		var image entity.Image
		image.Small = ImageProxy(gjson.Get(data, "image_urls.square_medium").String())
		image.Medium = ImageProxy(gjson.Get(data, "image_urls.medium").String())
		image.Large = ImageProxy(gjson.Get(data, "image_urls.large").String())
		image.Original = ImageProxy(gjson.Get(data, "meta_single_page.original_image_url").String())

		images = append(images, image)
	} else {
		g := GetInnerJSON(data, "meta_pages.#.image_urls")
		println(g)

		err := json.Unmarshal([]byte(g), &images)
		if err != nil {
			panic(err)
		}
	}

	return images
}

func ParseIllust(data string) entity.Illust {
	var illust entity.Illust

	images := ParseImages(data)
	illust.Images = images

	err := json.Unmarshal([]byte(data), &illust)
	if err != nil {
		panic("Failed to parse JSON")
	}

	return illust
}

func ParseIllusts(data string) []entity.Illust {
	var illusts []entity.Illust
	g := gjson.Get(data, "illusts").Array()

	for _, illust := range g {
		println(illust.String())
		illust := ParseIllust(illust.String())

		illusts = append(illusts, illust)
	}

	return illusts
}

func GetRecommendedIllust(c *gin.Context) []entity.Illust {
	URL := "https://hibi.cocomi.cf/api/pixiv/illust_recommended"

	s := Request(URL)
	illusts := ParseIllusts(s)

	return illusts
}

func GetRankingIllust(c *gin.Context, mode string) []entity.Illust {
	URL := "https://hibi.cocomi.cf/api/pixiv/rank?page=1&mode=" + mode

	s := Request(URL)
	illusts := ParseIllusts(s)

	return illusts
}

func GetNewestIllust(c *gin.Context) []entity.Illust {
	URL := "https://hibi.cocomi.cf/api/pixiv/illust_new"

	s := Request(URL)
	illusts := ParseIllusts(s)

	return illusts
}

func GetMemberIllust(c *gin.Context, id string) []entity.Illust {
	URL := "https://hibi.cocomi.cf/api/pixiv/member_illust?id=" + id

	s := Request(URL)
	illusts := ParseIllusts(s)

	return illusts
}

func GetRelatedIllust(c *gin.Context) []entity.Illust {
	id := c.Param("id")
	URL := "https://hibi.cocomi.cf/api/pixiv/related?id=" + id

	s := Request(URL)
	illusts := ParseIllusts(s)

	return illusts
}

func GetIllustByID(c *gin.Context) entity.Illust {
	id := c.Param("id")
	URL := "https://hibi.cocomi.cf/api/pixiv/illust?id=" + id

	s := Request(URL)
	g := GetInnerJSON(s, "illust")
	illust := ParseIllust(g)

	return illust
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
