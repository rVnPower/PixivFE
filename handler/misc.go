package handler

import (
	"errors"
	"fmt"
	"pixivfe/models"
	"strings"

	"github.com/goccy/go-json"
)

func (p *PixivClient) GetNewestArtworks(worktype string, r18 string) ([]models.IllustShort, error) {
	var pr models.PixivResponse
	var newWorks []models.IllustShort
	lastID := "0"

	for i := 0; i < 10; i++ {
		url := fmt.Sprintf(ArtworkNewestURL, worktype, r18, lastID)

		s, err := p.TextRequest(url)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(s), &pr)
		if err != nil {
			return nil, err
		}

		var body struct {
			Illusts []models.IllustShort `json:"illusts"`
			LastID  string               `json:"lastId"`
		}

		err = json.Unmarshal([]byte(pr.Body), &body)
		if err != nil {
			return nil, err
		}
		newWorks = append(newWorks, body.Illusts...)

		lastID = body.LastID
	}

	return newWorks, nil
}

func (p *PixivClient) GetRanking(mode string, content string, page string) (models.RankingResponse, error) {
	// Ranking data is formatted differently
	var pr models.RankingResponse

	url := fmt.Sprintf(ArtworkRankingURL, mode, content, page)

	s, err := p.TextRequest(url)

	if err != nil {
		return pr, err
	}

	err = json.Unmarshal([]byte(s), &pr)
	if err != nil {
		return pr, err
	}

	return pr, nil
}

func (p *PixivClient) GetSearch(artworkType string, name string, order string, age_settings string, page string) (*models.SearchResult, error) {
	var pr models.PixivResponse

	url := fmt.Sprintf(SearchArtworksURL, artworkType, name, order, age_settings, page)

	s, err := p.TextRequest(url)

	err = json.Unmarshal([]byte(s), &pr)

	if err != nil {
		return nil, err
	}

	// IDK how to do better than this lol
	temp := strings.ReplaceAll(string(pr.Body), `"illust"`, `"works"`)
	temp = strings.ReplaceAll(temp, `"manga"`, `"works"`)
	temp = strings.ReplaceAll(temp, `"illustManga"`, `"works"`)

	var resultRaw struct {
		*models.SearchResult
		ArtworksRaw json.RawMessage `json:"works"`
	}
	var artworks models.SearchArtworks
	var result *models.SearchResult

	err = json.Unmarshal([]byte(temp), &resultRaw)
	if err != nil {
		return nil, err
	}

	result = resultRaw.SearchResult

	err = json.Unmarshal([]byte(resultRaw.ArtworksRaw), &artworks)
	if err != nil {
		return nil, err
	}

	result.Artworks = artworks

	return result, nil
}

func (p *PixivClient) GetDiscoveryArtwork(mode string, count int) ([]models.IllustShort, error) {
	var artworks []models.IllustShort

	for count > 0 {
		var pr models.PixivResponse
		itemsForRequest := Min(100, count)

		count -= itemsForRequest

		url := fmt.Sprintf(ArtworkDiscoveryURL, mode, itemsForRequest)
		s, err := p.TextRequest(url)

		if err != nil {
			return artworks, err
		}

		err = json.Unmarshal([]byte(s), &pr)

		if pr.Error {
			return artworks, errors.New(pr.Message)
		}

		var thumbnail struct {
			Data json.RawMessage `json:"thumbnails"`
		}

		err = json.Unmarshal([]byte(pr.Body), &thumbnail)
		if err != nil {
			return nil, err
		}

		var body struct {
			Artworks []models.IllustShort `json:"illust"`
		}

		err = json.Unmarshal([]byte(thumbnail.Data), &body)
		if err != nil {
			return nil, err
		}

		artworks = append(artworks, body.Artworks...)
	}

	return artworks, nil
}
