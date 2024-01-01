package core

import (
	session "codeberg.org/vnpower/pixivfe/v2/core/config"
	http "codeberg.org/vnpower/pixivfe/v2/core/http"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/tidwall/gjson"
)

func GetDiscoveryArtwork(c *fiber.Ctx, mode string) ([]ArtworkBrief, error) {
	imageProxy := session.GetImageProxy(c)
	URL := http.GetDiscoveryURL(mode, 100)

	var artworks []ArtworkBrief

	resp, err := http.UnwrapWebAPIRequest(URL)
	println("here")
	if err != nil {
		return nil, err
	}

	resp = ProxyImages(resp, imageProxy)
	data := gjson.Get(resp, "thumbnails.illust").String()
	
	
	err = json.Unmarshal([]byte(data), &artworks)
	if err != nil {
		return nil, err
	}

	return artworks, nil
}
