package handler

import (
	"errors"
	"fmt"
	"math"
	"pixivfe/models"
	"sort"
	"strconv"

	"github.com/goccy/go-json"
)

func (p *PixivClient) GetUserArtworksID(id string, category string, page int) (string, int, error) {
	s, _ := p.TextRequest(fmt.Sprintf(UserArtworksURL, id))

	var pr models.PixivResponse

	err := json.Unmarshal([]byte(s), &pr)

	if err != nil {
		return "", -1, err
	}
	if pr.Error {
		return "", -1, errors.New(fmt.Sprintf("Pixiv returned error message: %s", pr.Message))
	}

	var ids []int
	var idsString string
	var body struct {
		Illusts json.RawMessage `json:"illusts"`
		Mangas  json.RawMessage `json:"manga"`
	}

	err = json.Unmarshal(pr.Body, &body)
	if err != nil {
		return "", -1, err
	}

	var illusts map[int]string
	var mangas map[int]string
	count := 0

	if err = json.Unmarshal(body.Illusts, &illusts); err != nil {
		illusts = make(map[int]string)
	}
	if err = json.Unmarshal(body.Mangas, &mangas); err != nil {
		mangas = make(map[int]string)
	}

	// Get the keys, because Pixiv only returns IDs (very evil)
	if category == "illustrations" || category == "artworks" {
		for k := range illusts {
			ids = append(ids, k)
			count++
		}
	}
	if category == "manga" || category == "artworks" {
		for k := range mangas {
			ids = append(ids, k)
			count++
		}
	}

	// Reverse sort the ids
	sort.Sort(sort.Reverse(sort.IntSlice(ids)))

	worksNumber := float64(count)
	worksPerPage := 30.0

	if page < 1 || float64(page) > math.Ceil(worksNumber/worksPerPage)+1.0 {
		return "", -1, errors.New("Page overflow")
	}

	start := (page - 1) * int(worksPerPage)
	end := int(math.Min(float64(page)*worksPerPage, worksNumber)) // no overflow

	for _, k := range ids[start:end] {
		idsString += fmt.Sprintf("&ids[]=%d", k)
	}

	return idsString, count, nil
}

func (p *PixivClient) GetUserArtworks(id string, ids string) ([]models.IllustShort, error) {
	url := fmt.Sprintf(UserArtworksFullURL, id, ids)

	var pr models.PixivResponse
	var works []models.IllustShort

	s, err := p.TextRequest(url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(s), &pr)
	if err != nil {
		return nil, err
	}

	var body struct {
		Illusts map[int]json.RawMessage `json:"works"`
	}

	err = json.Unmarshal(pr.Body, &body)
	if err != nil {
		return nil, err
	}

	for _, v := range body.Illusts {
		var illust models.IllustShort
		err = json.Unmarshal(v, &illust)

		works = append(works, illust)
	}

	return works, nil
}

func (p *PixivClient) GetUserBasicInformation(id string) (models.UserShort, error) {
	var pr models.PixivResponse
	var user models.UserShort

	s, _ := p.TextRequest(fmt.Sprintf(UserBasicInformationURL, id))

	err := json.Unmarshal([]byte(s), &pr)
	if err != nil {
		return user, err
	}
	if pr.Error {
		return user, errors.New(fmt.Sprintf("Pixiv returned error message: %s", pr.Message))
	}

	err = json.Unmarshal([]byte(pr.Body), &user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (p *PixivClient) GetUserInformation(id string, category string, page int) (*models.User, error) {
	var user *models.User
	var pr models.PixivResponse

	ids, count, err := p.GetUserArtworksID(id, category, page)
	if err != nil {
		return nil, err
	}

	s, _ := p.TextRequest(fmt.Sprintf(UserInformationURL, id))

	err = json.Unmarshal([]byte(s), &pr)

	if err != nil {
		return nil, err
	}
	if pr.Error {
		return nil, errors.New(fmt.Sprintf("Pixiv returned error message: %s", pr.Message))
	}

	var body struct {
		*models.User
		Background map[string]interface{} `json:"background"`
	}

	// Basic user information
	err = json.Unmarshal([]byte(pr.Body), &body)
	if err != nil {
		return nil, err
	}

	user = body.User

	// Artworks
	works, _ := p.GetUserArtworks(id, ids)
	// IDK but the order got shuffled even though Pixiv sorted the IDs in the response
	sort.Slice(works[:], func(i, j int) bool {
		left, _ := strconv.Atoi(works[i].ID)
		right, _ := strconv.Atoi(works[j].ID)
		return left > right
	})
	user.Artworks = works

	// Background image
	if body.Background != nil {
		user.BackgroundImage = body.Background["url"].(string)
	}

	// Artworks count
	user.ArtworksCount = count

	// Frequent tags
	user.FrequentTags, err = p.GetFrequentTags(ids)

	return user, nil
}
