package core

import (
	"fmt"

	session "codeberg.org/vnpower/pixivfe/v2/core/user"
	http "codeberg.org/vnpower/pixivfe/v2/core/http"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/tidwall/gjson"
)

func GetDiscoveryArtwork(c *fiber.Ctx, mode string) ([]ArtworkBrief, error) {
	token := session.GetCookie(c, session.Cookie_ImageProxy)

	URL := http.GetDiscoveryURL(mode, 100)

	var artworks []ArtworkBrief

	resp, err := http.UnwrapWebAPIRequest(c.Context(), URL, token)
	if err != nil {
		return nil, err
	}
	resp = session.ProxyImageUrl(c, resp)
	if !gjson.Valid(resp) {
		return nil, fmt.Errorf("Invalid JSON: %v", resp)
	}
	data := gjson.Get(resp, "thumbnails.illust").String()

	err = json.Unmarshal([]byte(data), &artworks)
	if err != nil {
		return nil, err
	}

	return artworks, nil
}

func GetDiscoveryNovels(c *fiber.Ctx, mode string) ([]NovelBrief, error) {
	token := session.GetCookie(c, session.Cookie_ImageProxy)

	URL := http.GetDiscoveryNovelURL(mode, 100)

	var novels []NovelBrief

	resp, err := http.UnwrapWebAPIRequest(c.Context(), URL, token)
	if err != nil {
		return nil, err
	}
	resp = session.ProxyImageUrl(c, resp)
	if !gjson.Valid(resp) {
		return nil, fmt.Errorf("Invalid JSON: %v", resp)
	}
	data := gjson.Get(resp, "thumbnails.novel").String()

	err = json.Unmarshal([]byte(data), &novels)
	if err != nil {
		return nil, err
	}

	return novels, nil
}
