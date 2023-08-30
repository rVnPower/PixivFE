package handler

import (
	"fmt"
	"sort"
	"strconv"

	"codeberg.org/vnpower/pixivfe/models"
	"github.com/goccy/go-json"
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

	// Begin testing here

	c1 := make(chan []models.Image)
	c2 := make(chan []models.IllustShort)
	c3 := make(chan models.UserShort)
	c4 := make(chan []models.Tag)
	c5 := make(chan []models.IllustShort)
	c6 := make(chan []models.Comment)

	go func() {
		// Get illust images
		images, err = p.GetArtworkImages(id)
		if err != nil {
			c1 <- nil
		}
		c1 <- images
	}()

	go func() {
		// Get recent artworks
		ids := make([]int, 0)

		for k := range illust.Recent {
			ids = append(ids, k)
		}

		sort.Sort(sort.Reverse(sort.IntSlice(ids)))

		idsString := ""
		count := min(len(ids), 20)

		for i := 0; i < count; i++ {
			idsString += fmt.Sprintf("&ids[]=%d", ids[i])
		}

		recent, err := p.GetUserArtworks(illust.UserID, idsString)
		if err != nil {
			c2 <- nil
		}
		sort.Slice(recent[:], func(i, j int) bool {
			left, _ := strconv.Atoi(recent[i].ID)
			right, _ := strconv.Atoi(recent[j].ID)
			return left > right
		})
		c2 <- recent

	}()

	go func() {
		// Get basic user information (the URL above does not contain avatars)
		userInfo, err := p.GetUserBasicInformation(illust.UserID)
		if err != nil {
			//
		}
		c3 <- userInfo
	}()

	go func() {
		var tagsList []models.Tag
		// Extract tags
		var tags struct {
			Tags []struct {
				Tag         string            `json:"tag"`
				Translation map[string]string `json:"translation"`
			} `json:"tags"`
		}
		err = json.Unmarshal(illust.RawTags, &tags)
		if err != nil {
			c4 <- nil
		}

		for _, tag := range tags.Tags {
			var newTag models.Tag
			newTag.Name = tag.Tag
			newTag.TranslatedName = tag.Translation["en"]

			tagsList = append(tagsList, newTag)
		}
		c4 <- tagsList
	}()

	go func() {
		related, _ := p.GetRelatedArtworks(id)
		// Error handling...
		c5 <- related
	}()

	go func() {
		comments, _ := p.GetArtworkComments(id)
		// Error handling...
		c6 <- comments
	}()

	illust.Images = <-c1
	illust.RecentWorks = <-c2
	illust.User = <-c3
	illust.Tags = <-c4
	illust.RelatedWorks = <-c5
	illust.CommentsList = <-c6

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
	if err != nil {
		return nil, err
	}

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
