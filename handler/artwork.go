package handler

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/goccy/go-json"
	"pixivfe/models"
)

func (p *PixivClient) GetArtworkImages(id string) ([]models.Image, error) {
	var resp []models.ImageResponse
	var images []models.Image

	URL := fmt.Sprintf(ArtworkImagesURL, id)

	response, err := p.PixivRequest(URL)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(response), &resp)
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
	var images []models.Image

	URL := fmt.Sprintf(ArtworkInformationURL, id)

	response, err := p.PixivRequest(URL)
	if err != nil {
		return nil, err
	}

	var illust struct {
		*models.Illust

		Recent  map[int]any     `json:"userIllusts"`
		RawTags json.RawMessage `json:"tags"`
	}

	// Parse basic illust information
	err = json.Unmarshal([]byte(response), &illust)
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
	ids := make([]int, len(illust.Recent))

	for k := range illust.Recent {
		ids = append(ids, k)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(ids)))

	idsString := ""
	count := Min(len(ids), 30)

	for i := 0; i < count; i++ {
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
	var body struct {
		Comments []models.Comment `json:"comments"`
	}

	URL := fmt.Sprintf(ArtworkCommentsURL, id)

	response, err := p.PixivRequest(URL)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(response), &body)

	return body.Comments, nil
}

func (p *PixivClient) GetRelatedArtworks(id string) ([]models.IllustShort, error) {
	var body struct {
		Illusts []models.IllustShort `json:"illusts"`
	}

	URL := fmt.Sprintf(ArtworkRelatedURL, id, 96)

	response, err := p.PixivRequest(URL)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(response), &body)
	if err != nil {
		return nil, err
	}

	return body.Illusts, nil
}
