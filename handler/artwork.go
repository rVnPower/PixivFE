package handler

import (
	"errors"
	"fmt"
	"sort"
	"strconv"

	"github.com/goccy/go-json"
	"pixivfe/models"
)

func (p *PixivClient) GetArtworkImages(id string) ([]models.Image, error) {
	s, _ := p.TextRequest(fmt.Sprintf(ArtworkImagesURL, id))

	var pr models.PixivResponse
	var resp []models.ImageResponse
	var images []models.Image

	err := json.Unmarshal([]byte(s), &pr)
	if err != nil {
		return images, err
	}
	if pr.Error {
		return images, errors.New(fmt.Sprintf("Pixiv returned error message: %s", pr.Message))
	}

	err = json.Unmarshal([]byte(pr.Body), &resp)
	if err != nil {
		return images, err
	}

	// Extract and proxy every images
	for _, imageRaw := range resp {
		var image models.Image

		image.Small = imageRaw.Urls["thumb_mini"]
		image.Medium = imageRaw.Urls["small"]
		image.Large = imageRaw.Urls["regular"]
		image.Original = imageRaw.Urls["original"]

		images = append(images, image)
	}

	return images, nil
}

func (p *PixivClient) GetArtworkByID(id string) (*models.Illust, error) {
	s, _ := p.TextRequest(fmt.Sprintf(ArtworkInformationURL, id))

	var pr models.PixivResponse
	var images []models.Image

	// Parse Pixiv response body
	err := json.Unmarshal([]byte(s), &pr)
	if err != nil {
		return nil, err
	}
	if pr.Error {
		return nil, errors.New(fmt.Sprintf("Pixiv returned error message: %s", pr.Message))
	}

	var illust struct {
		*models.Illust

		Recent  map[int]any     `json:"userIllusts"`
		RawTags json.RawMessage `json:"tags"`
	}

	// Parse basic illust information
	err = json.Unmarshal([]byte(pr.Body), &illust)
	if err != nil {
		return nil, err
	}

	// Get illust images
	images, err = p.GetArtworkImages(id)
	if err != nil {
		return nil, err
	}

	illust.Images = images

	// Get recent artworks
	var ids []int
	idsString := ""

	for k := range illust.Recent {
		ids = append(ids, k)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(ids)))

	count := len(ids)
	for i := 0; i < 30 && i < count; i++ {
		idsString += fmt.Sprintf("&ids[]=%d", ids[i])
	}

	recent, err := p.GetUserArtworks(illust.UserID, idsString)
	if err != nil {
		return nil, err
	}
	sort.Slice(recent[:], func(i, j int) bool {
		left, _ := strconv.Atoi(recent[i].ID)
		right, _ := strconv.Atoi(recent[j].ID)
		return left > right
	})

	illust.RecentWorks = recent

	// Get basic user information (the URL above does not contain avatars)
	userInfo, err := p.GetUserBasicInformation(illust.UserID)
	if err != nil {
		return nil, err
	}

	illust.User = userInfo

	// Extract tags
	var tags struct {
		Tags []struct {
			Tag         string            `json:"tag"`
			Translation map[string]string `json:"translation"`
		} `json:"tags"`
	}
	err = json.Unmarshal(illust.RawTags, &tags)
	if err != nil {
		return nil, err
	}

	for _, tag := range tags.Tags {
		var newTag models.Tag
		newTag.Name = tag.Tag
		newTag.TranslatedName = tag.Translation["en"]

		illust.Tags = append(illust.Tags, newTag)
	}

	return illust.Illust, nil
}

func (p *PixivClient) GetArtworkComments(id string) ([]models.Comment, error) {
	var pr models.PixivResponse
	var body struct {
		Comments []models.Comment `json:"comments"`
	}

	s, _ := p.TextRequest(fmt.Sprintf(ArtworkCommentsURL, id))

	err := json.Unmarshal([]byte(s), &pr)

	if err != nil {
		return nil, err
	}
	if pr.Error {
		return nil, errors.New(fmt.Sprintf("Pixiv returned error message: %s", pr.Message))
	}

	err = json.Unmarshal([]byte(pr.Body), &body)

	return body.Comments, nil
}

func (p *PixivClient) GetRelatedArtworks(id string) ([]models.IllustShort, error) {
	url := fmt.Sprintf(ArtworkRelatedURL, id, 30)

	var pr models.PixivResponse
	var body struct {
		Illusts []models.IllustShort `json:"illusts"`
	}

	s, _ := p.TextRequest(url)

	err := json.Unmarshal([]byte(s), &pr)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(pr.Body), &body)
	if err != nil {
		return nil, err
	}

	return body.Illusts, nil
}
