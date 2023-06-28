package handler

import (
	"fmt"
	"pixivfe/models"

	"github.com/goccy/go-json"
)

func (p *PixivClient) GetLandingPage(mode string) (models.LandingArtworks, error) {
	var context models.LandingArtworks
	URL := fmt.Sprintf(LandingPageURL, mode)

	response, err := p.PixivRequest(URL)
	if err != nil {
		return context, err
	}

	var pages struct {
		Follow     []any `json:"follow"`
		Commission []any `json:"completeRequestIds"`
	}

	var body struct {
		Thumbnails json.RawMessage `json:"thumbnails"`
		Page       json.RawMessage `json:"page"`
	}

	var artworks struct {
		Artworks []models.IllustShort `json:"illust"`
	}

	err = json.Unmarshal([]byte(response), &body)
	if err != nil {
		return context, err
	}
	err = json.Unmarshal([]byte(body.Thumbnails), &artworks)
	if err != nil {
		return context, err
	}
	err = json.Unmarshal([]byte(body.Page), &pages)
	if err != nil {
		return context, err
	}

	// Keep track
	count := len(pages.Commission)

	context.Commissions = artworks.Artworks[:count]
	context.Following = artworks.Artworks[count:len(pages.Follow)]

	count += len(pages.Follow)

	return context, nil
}
