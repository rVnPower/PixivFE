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

		err := json.Unmarshal([]byte(g), &images)
		if err != nil {
			panic(err)
		}
	}

	return images
}

func ParseIllust(data string) []entity.Illust {
	var illusts []entity.Illust
	g := gjson.Get(data, "#")

	for _, handle := range g.Array() {
		println(handle.String())
		println("here")
	}
	// g.ForEach(func(key, value gjson.Result) bool {
	// 	var illust entity.Illust
	// 	v := value.String()
	// 	println(v)
	// 	err := json.Unmarshal([]byte(v), &illust)
	//
	// 	if err != nil {
	// 		panic(err)
	// 	}
	//
	// 	illust.Images = ParseImages(v)
	// 	illusts = append(illusts, illust)
	// 	return true
	// })
	//
	// println(illusts)

	return illusts
}

func GetRecommendedIllust(c *gin.Context) []entity.Illust {
	URL := "https://hibi.cocomi.cf/api/pixiv/illust_recommended"
	var illusts []entity.Illust

	s := Request(URL)
	ParseIllust(s)
	g := GetInnerJSON(s, "illusts")

	err := json.Unmarshal([]byte(g), &illusts)
	if err != nil {
		panic("Failed to parse JSON")
	}

	return illusts
}

func GetRankingIllust(c *gin.Context, mode string) []entity.Illust {
	URL := "https://hibi.cocomi.cf/api/pixiv/rank?page=1&mode=" + mode
	var illusts []entity.Illust

	s := Request(URL)
	g := GetInnerJSON(s, "illusts")

	err := json.Unmarshal([]byte(g), &illusts)
	if err != nil {
		panic("Failed to parse JSON")
	}

	return illusts
}

func GetNewestIllust(c *gin.Context) []entity.Illust {
	URL := "https://hibi.cocomi.cf/api/pixiv/illust_new"
	var illusts []entity.Illust

	s := Request(URL)
	g := GetInnerJSON(s, "illusts")

	err := json.Unmarshal([]byte(g), &illusts)
	if err != nil {
		panic("Failed to parse JSON")
	}

	return illusts
}

func GetMemberIllust(c *gin.Context, id string) []entity.Illust {
	URL := "https://hibi.cocomi.cf/api/pixiv/member_illust?id=" + id
	var illusts []entity.Illust

	s := Request(URL)
	g := GetInnerJSON(s, "illusts")

	err := json.Unmarshal([]byte(g), &illusts)
	if err != nil {
		panic("Failed to parse JSON")
	}

	return illusts
}

func GetRelatedIllust(c *gin.Context) []entity.Illust {
	id := c.Param("id")
	URL := "https://hibi.cocomi.cf/api/pixiv/related?id=" + id
	var illusts []entity.Illust

	s := Request(URL)
	g := GetInnerJSON(s, "illusts")

	err := json.Unmarshal([]byte(g), &illusts)
	if err != nil {
		panic("Failed to parse JSON")
	}

	return illusts
}

func GetIllustByID(c *gin.Context) entity.Illust {
	id := c.Param("id")
	URL := "https://hibi.cocomi.cf/api/pixiv/illust?id=" + id
	var illust entity.Illust

	s := Request(URL)
	g := GetInnerJSON(s, "illust")

	err := json.Unmarshal([]byte(g), &illust)
	if err != nil {
		panic(err)
	}

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
