package handler

import (
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"pixivfe/models"
)

func (p *PixivClient) GetTagData(name string) (models.TagDetail, error) {
	var pr models.PixivResponse
	var tag models.TagDetail

	url := fmt.Sprintf(SearchTagURL, name)

	s, err := p.TextRequest(url)

	err = json.Unmarshal([]byte(s), &pr)

	if err != nil {
		return tag, err
	}

	if pr.Error {
		return tag, errors.New(fmt.Sprintf("Pixiv returned error message: %s", pr.Message))
	}

	err = json.Unmarshal([]byte(pr.Body), &tag)
	if err != nil {
		return tag, err
	}

	return tag, nil
}

func (p *PixivClient) GetFrequentTags(ids string) ([]models.FrequentTag, error) {
	s, _ := p.TextRequest(fmt.Sprintf(FrequentTagsURL, ids))

	var pr models.PixivResponse
	var tags []models.FrequentTag
	// Parse Pixiv response body
	err := json.Unmarshal([]byte(s), &pr)
	if err != nil {
		return nil, err
	}
	if pr.Error {
		return nil, errors.New(fmt.Sprintf("Pixiv returned error message: %s", pr.Message))
	}

	err = json.Unmarshal([]byte(pr.Body), &tags)
	if err != nil {
		return nil, err
	}

	return tags, nil
}
