package handler

import (
	"fmt"

	"codeberg.org/vnpower/pixivfe/models"
	"github.com/goccy/go-json"
)

func (p *PixivClient) GetNewestFromFollowing(mode, page, token string) ([]models.IllustShort, error) {
	URL := fmt.Sprintf(NewestFromFollowURL, "illust", mode, page)

	var body struct {
		Thumbnails json.RawMessage `json:"thumbnails"`
	}

	var artworks struct {
		Artworks []models.IllustShort `json:"illust"`
	}

	response, err := p.PixivRequest(URL, token)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(response), &body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(body.Thumbnails), &artworks)
	if err != nil {
		return nil, err
	}

	return artworks.Artworks, nil
}

// func (p *PixivClient) FollowUser(id string) error {
// 	formData := url.Values{}
// 	formData.Add("mode", "add")
// 	formData.Add("type", "user")
// 	formData.Add("user_id", id)
// 	formData.Add("tag", "")
// 	formData.Add("restrict", "0")
// 	formData.Add("format", "json")

// 	init, err := p.GetCSRF()
// 	println(init)
// 	if err != nil {
// 		return err
// 	}

// 	pattern := regexp.MustCompile(`.*pixiv.context.token = "([a-z0-9]{32})"?.*`)
// 	quotesPattern := regexp.MustCompile(`([a-z0-9]{32})`)
// 	token := quotesPattern.FindString(pattern.FindString(init))
// 	println(token)

// 	_, err = p.RequestWithFormData(FollowUserURL, formData, token)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
