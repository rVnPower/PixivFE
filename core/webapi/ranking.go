package core

import (
	"errors"
	"strings"

	session "codeberg.org/vnpower/pixivfe/v2/core/config"
	http "codeberg.org/vnpower/pixivfe/v2/core/http"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

type Ranking struct {
	Contents []struct {
		Title        string `json:"title"`
		Image        string `json:"url"`
		Pages        string `json:"illust_page_count"`
		ArtistName   string `json:"user_name"`
		ArtistAvatar string `json:"profile_img"`
		ID           int    `json:"illust_id"`
		ArtistID     int    `json:"user_id"`
		Rank         int    `json:"rank"`
	} `json:"contents"`

	Mode        string          `json:"mode"`
	Content     string          `json:"content"`
	Page        int             `json:"page"`
	RankTotal   int             `json:"rank_total"`
	CurrentDate string          `json:"date"`
	PrevDateRaw json.RawMessage `json:"prev_date"`
	NextDateRaw json.RawMessage `json:"next_date"`
	PrevDate    string
	NextDate    string
}

func GetRanking(c *fiber.Ctx, mode, content, date, page string) (Ranking, error) {
	URL := http.GetRankingURL(mode, content, date, page)

	var ranking Ranking

	resp := http.WebAPIRequest(URL, "")
	if !resp.Ok {
		return ranking, errors.New(resp.Message)
	}
	proxiedResp := session.ProxyImageUrl(resp.Body)

	err := json.Unmarshal([]byte(proxiedResp), &ranking)
	if err != nil {
		return ranking, err
	}

	ranking.PrevDate = strings.ReplaceAll(string(ranking.PrevDateRaw[:]), "\"", "")
	ranking.NextDate = strings.ReplaceAll(string(ranking.NextDateRaw[:]), "\"", "")

	return ranking, nil
}
