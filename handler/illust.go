package handler

import (
	"encoding/json"
	"net/http"
	"pixivfe/entity"

	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

func GetRecommendedIllust(c *gin.Context) []entity.Illust {
	URL := "https://hibi.cocomi.cf/api/pixiv/illust_recommended"
	var illusts []entity.Illust

	s := Request(URL)
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
