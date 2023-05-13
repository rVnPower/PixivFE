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

func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
