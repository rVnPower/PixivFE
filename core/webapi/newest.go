package core

import (
	"html/template"
	"time"
	"strings"

	session "codeberg.org/vnpower/pixivfe/v2/core/config"
	http "codeberg.org/vnpower/pixivfe/v2/core/http"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

type ArtworkBrief struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Description  template.HTML `json:"description"`
	ArtistID     string        `json:"userId"`
	ArtistName   string        `json:"userName"`
	ArtistAvatar string        `json:"profileImageUrl"`
	Date         time.Time     `json:"uploadDate"`
	Thumbnail    string        `json:"url"`
	Pages        int           `json:"pageCount"`
	XRestrict    int           `json:"xRestrict"`
	AiType       int           `json:"aiType"`
}

func ProxyImages(s, proxy string) string {
	proxy = "https://" + proxy
	s = strings.ReplaceAll(s, `https:\/\/i.pximg.net`, proxy)
	s = strings.ReplaceAll(s, `https:\/\/s.pximg.net`, proxy)

	return s
}

func GetNewestArtworks(c *fiber.Ctx, worktype string, r18 string) ([]ArtworkBrief, error) {
	imageProxy := session.GetImageProxy(c)
	URL := http.GetNewestArtworksURL(worktype, r18, "0")

	var body struct {
		Artworks []ArtworkBrief `json:"illusts"`
		// LastId string
	}

	resp, err := http.UnwrapWebAPIRequest(URL)
	if err != nil {
		return nil, err
	}

	resp = ProxyImages(resp, imageProxy)

	err = json.Unmarshal([]byte(resp), &body)
	if err != nil {
		return nil, err
	}

	return body.Artworks, nil
}
