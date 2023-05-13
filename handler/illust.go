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
	resp, err := http.Get(URL)

	if err != nil {
		panic(":(")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(":(")
	}
	s := string(body)
	g := gjson.Get(s, "illusts").String()

	err = json.Unmarshal([]byte(g), &illusts)
	if err != nil {
		panic(":(")
	}

	return illusts
}

func GetRankingIllust(c *gin.Context, mode string) []entity.Illust {
	if mode == "" {
		mode = "day"
	}
	URL := "https://hibi.cocomi.cf/api/pixiv/rank?mode=" + mode
	var illusts []entity.Illust
	illusts = append(illusts, entity.Illust{}) // Placeholder to shift the index

	resp, err := http.Get(URL)

	if err != nil {
		panic("Failed to request")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic("Failed to parse body")
	}
	s := string(body)
	g := gjson.Get(s, "illusts").String()

	err = json.Unmarshal([]byte(g), &illusts)
	if err != nil {
		panic("Failed to parse Json.")
	}

	return illusts
}
