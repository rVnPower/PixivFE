package handler

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"

	"codeberg.org/vnpower/pixivfe/models"
	"github.com/goccy/go-json"
)

func (p *PixivClient) GetUserArtworksID(id string, category string, page int) (string, int, error) {
	URL := fmt.Sprintf(UserArtworksURL, id)

	response, err := p.PixivRequest(URL)
	if err != nil {
		return "", -1, err
	}

	var ids []int
	var idsString string
	var body struct {
		Illusts json.RawMessage `json:"illusts"`
		Mangas  json.RawMessage `json:"manga"`
	}

	err = json.Unmarshal(response, &body)
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
	end := int(min(float64(page)*worksPerPage, worksNumber)) // no overflow

	for _, k := range ids[start:end] {
		idsString += fmt.Sprintf("&ids[]=%d", k)
	}

	return idsString, count, nil
}

func (p *PixivClient) GetUserArtworks(id string, ids string) ([]models.IllustShort, error) {
	var works []models.IllustShort

	URL := fmt.Sprintf(UserArtworksFullURL, id, ids)

	response, err := p.PixivRequest(URL)
	if err != nil {
		return nil, err
	}

	var body struct {
		Illusts map[int]json.RawMessage `json:"works"`
	}

	err = json.Unmarshal(response, &body)
	if err != nil {
		return nil, err
	}

	for _, v := range body.Illusts {
		var illust models.IllustShort
		err = json.Unmarshal(v, &illust)

		if err != nil {
			return nil, err
		}

		works = append(works, illust)
	}

	return works, nil
}

func (p *PixivClient) GetUserBasicInformation(id string) (models.UserShort, error) {
	var user models.UserShort

	URL := fmt.Sprintf(UserBasicInformationURL, id)

	response, err := p.PixivRequest(URL)
	if err != nil {
		return user, err
	}

	err = json.Unmarshal([]byte(response), &user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (p *PixivClient) GetUserInformation(id string, category string, page int) (*models.User, error) {
	var user *models.User

	URL := fmt.Sprintf(UserInformationURL, id)

	response, err := p.PixivRequest(URL)
	if err != nil {
		return user, err
	}

	var body struct {
		*models.User
		Background map[string]interface{} `json:"background"`
	}

	// Basic user information
	err = json.Unmarshal([]byte(response), &body)
	if err != nil {
		return nil, err
	}

	user = body.User


	if category != "bookmarks" {
		// Artworks
		ids, count, err := p.GetUserArtworksID(id, category, page)
		if err != nil {
			return nil, err
		}
		works, _ := p.GetUserArtworks(id, ids)
		// IDK but the order got shuffled even though Pixiv sorted the IDs in the response
		sort.Slice(works[:], func(i, j int) bool {
			left, _ := strconv.Atoi(works[i].ID)
			right, _ := strconv.Atoi(works[j].ID)
			return left > right
		})
		user.Artworks = works

		// Artworks count
		user.ArtworksCount = count

		// Frequent tags
		user.FrequentTags, err = p.GetFrequentTags(ids)

		if err != nil {
			return nil, err
		}
	} else {
		// Bookmarks
		works, count, err := p.GetUserBookmarks(id, "show", page)
		if err != nil {
			return nil, err
		}

		user.Artworks = works

		// Public bookmarks count
		user.ArtworksCount = count

		// Parse social medias
		user.ParseSocial()
	}

	// Background image
	if body.Background != nil {
		user.BackgroundImage = body.Background["url"].(string)
	}

	return user, nil
}

func (p *PixivClient) GetUserBookmarks(id string, mode string, page int) ([]models.IllustShort, int, error) {
	page--
	URL := fmt.Sprintf(UserBookmarksURL, id, page*48, mode)

	response, err := p.PixivRequest(URL)
	if err != nil {
		return nil, -1, err
	}

	var body struct {
		Artworks []json.RawMessage `json:"works"`
		Total    int               `json:"total"`
	}

	err = json.Unmarshal([]byte(response), &body)
	if err != nil {
		return nil, -1, err
	}

	artworks := make([]models.IllustShort, len(body.Artworks))

	for index, value := range body.Artworks {
		var artwork models.IllustShort

		err = json.Unmarshal([]byte(value), &artwork)
		if err != nil {
			artworks[index] = models.IllustShort{
				ID:        "#",
				Title:     "Deleted or Private",
				Thumbnail: "https://s.pximg.net/common/images/limit_unknown_360.png",
			}
			continue
		}
		artworks[index] = artwork
	}

	return artworks, body.Total, nil
}
