package core

import (
	session "codeberg.org/vnpower/pixivfe/v2/core/config"
	http "codeberg.org/vnpower/pixivfe/v2/core/http"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

func GetNewestFromFollowing(c *fiber.Ctx, mode, page string) ([]ArtworkBrief, error) {
	imageProxy := session.GetImageProxy(c)
	token := session.GetToken(c)
	URL := http.GetNewestFromFollowingURL(mode, page)

	var body struct {
		Thumbnails json.RawMessage `json:"thumbnails"`
	}

	var artworks struct {
		Artworks []ArtworkBrief `json:"illust"`
	}

	resp, err := http.UnwrapWebAPIRequest(URL, token)
	if err != nil {
		return nil, err
	}

	resp = ProxyImages(resp, imageProxy)

	err = json.Unmarshal([]byte(resp), &body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(body.Thumbnails), &artworks)
	if err != nil {
		return nil, err
	}

	return artworks.Artworks, nil
}
