package core

import (
	session "codeberg.org/vnpower/pixivfe/v2/core/user"
	http "codeberg.org/vnpower/pixivfe/v2/core/http"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

func GetNewestArtworks(c *fiber.Ctx, worktype string, r18 string) ([]ArtworkBrief, error) {
	token := session.GetToken(c)
	URL := http.GetNewestArtworksURL(worktype, r18, "0")

	var body struct {
		Artworks []ArtworkBrief `json:"illusts"`
		// LastId string
	}

	resp, err := http.UnwrapWebAPIRequest(c.Context(), URL, token)
	if err != nil {
		return nil, err
	}
	resp = session.ProxyImageUrl(c, resp)

	err = json.Unmarshal([]byte(resp), &body)
	if err != nil {
		return nil, err
	}

	return body.Artworks, nil
}
