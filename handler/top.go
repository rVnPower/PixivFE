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

	type IDS struct {
		Ids []any `json:"ids"`
	}

	var pages struct {
		Pixivision        []models.Pixivision `json:"pixivision"`
		Follow            []any               `json:"follow"`
		Commission        []any               `json:"completeRequestIds"`
		Newest            []any               `json:"newPost"`
		Recommended       IDS                 `json:"recommend"`
		EditorRecommended []any               `json:"editorRecommend"`
		UserRecommended   []any               `json:"recommendUser"`
		RecommendedByTags []struct {
			models.LandingRecommendByTags
			Ids []any `json:"ids"`
		} `json:"recommendByTag"`
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

	context.Pixivision = pages.Pixivision

	// Keep track
	count := len(pages.Commission)

	context.Commissions = artworks.Artworks[:count]
	context.Following = artworks.Artworks[count : count+len(pages.Follow)]

	count += len(pages.Follow)

	context.Recommended = artworks.Artworks[count : count+len(pages.Recommended.Ids)]
	count += len(pages.Recommended.Ids)

	context.Newest = artworks.Artworks[count : count+len(pages.Newest)]
	count += len(pages.Newest)

	// For rankings, we just take 100 anyway
	context.Rankings = artworks.Artworks[count : count+100]
	count += 100

	// IDK what this is
	count += len(pages.EditorRecommended)

	context.Users = artworks.Artworks[count : count+len(pages.UserRecommended)*3]

	count += len(pages.UserRecommended) * 3

	for i := 0; i < len(pages.RecommendedByTags); i++ {
		temp := pages.RecommendedByTags[i]
		temp.Artworks = artworks.Artworks[count : count+min(len(temp.Ids), 18)]

		context.RecommendByTags = append(context.RecommendByTags, temp.LandingRecommendByTags)
		count += len(temp.Ids)
	}

	return context, nil
}
