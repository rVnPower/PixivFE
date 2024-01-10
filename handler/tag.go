package handler

import (
	"codeberg.org/vnpower/pixivfe/models"
	"github.com/goccy/go-json"
)

func (p *PixivClient) GetTagData(name string) (models.TagDetail, error) {
	var tag models.TagDetail

	URL := UrlSprintf(SearchTagURL, name)

	response, err := p.PixivRequest(URL)
	if err != nil {
		return tag, err
	}

	err = json.Unmarshal([]byte(response), &tag)
	if err != nil {
		return tag, err
	}

	return tag, nil
}

func (p *PixivClient) GetFrequentTags(ids string) ([]models.FrequentTag, error) {
	var tags []models.FrequentTag

	URL := UrlSprintf(FrequentTagsURL, ids)

	response, err := p.PixivRequest(URL)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(response), &tags)
	if err != nil {
		return nil, err
	}

	return tags, nil
}
