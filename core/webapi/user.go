package core

import (
	"errors"
	"fmt"
	"html/template"
	"math"
	"sort"

	session "codeberg.org/vnpower/pixivfe/v2/core/session"
	http "codeberg.org/vnpower/pixivfe/v2/core/http"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

type FrequentTag struct {
	Name           string `json:"tag"`
	TranslatedName string `json:"tag_translation"`
}

type User struct {
	ID              string                 `json:"userId"`
	Name            string                 `json:"name"`
	Avatar          string                 `json:"imageBig"`
	Following       int                    `json:"following"`
	MyPixiv         int                    `json:"mypixivCount"`
	Comment         template.HTML          `json:"commentHtml"`
	Webpage         string                 `json:"webpage"`
	SocialRaw       json.RawMessage        `json:"social"`
	Artworks        []ArtworkBrief         `json:"artworks"`
	Background      map[string]interface{} `json:"background"`
	ArtworksCount   int
	FrequentTags    []FrequentTag
	Social          map[string]map[string]string
	BackgroundImage string
}

func (s *User) ParseSocial() error {
	if string(s.SocialRaw[:]) == "[]" {
		// Fuck Pixiv
		return nil
	}

	err := json.Unmarshal(s.SocialRaw, &s.Social)
	if err != nil {
		return err
	}
	return nil
}

func GetFrequentTags(c *fiber.Ctx, ids string) ([]FrequentTag, error) {
	var tags []FrequentTag

	URL := http.GetFrequentTagsURL(ids)

	response, err := http.UnwrapWebAPIRequest(c.Context(), URL, "")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(response), &tags)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func GetUserArtworks(c *fiber.Ctx, id, ids string) ([]ArtworkBrief, error) {
	var works []ArtworkBrief

	URL := http.GetUserFullArtworkURL(id, ids)

	resp, err := http.UnwrapWebAPIRequest(c.Context(), URL, "")
	if err != nil {
		return nil, err
	}
	resp = session.ProxyImageUrl(c, resp)

	var body struct {
		Illusts map[int]json.RawMessage `json:"works"`
	}

	err = json.Unmarshal([]byte(resp), &body)
	if err != nil {
		return nil, err
	}

	for _, v := range body.Illusts {
		var illust ArtworkBrief
		err = json.Unmarshal(v, &illust)

		if err != nil {
			return nil, err
		}

		works = append(works, illust)
	}

	return works, nil
}

func GetUserArtworksID(c *fiber.Ctx, id, category string, page int) (string, int, error) {
	URL := http.GetUserArtworksURL(id)

	resp, err := http.UnwrapWebAPIRequest(c.Context(), URL, "")
	if err != nil {
		return "", -1, err
	}

	var body struct {
		Illusts json.RawMessage `json:"illusts"`
		Mangas  json.RawMessage `json:"manga"`
	}

	err = json.Unmarshal([]byte(resp), &body)
	if err != nil {
		return "", -1, err
	}

	var ids []int
	var idsString string

	err = json.Unmarshal([]byte(resp), &body)
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
		return "", -1, errors.New("No page available.")
	}

	start := (page - 1) * int(worksPerPage)
	end := int(min(float64(page)*worksPerPage, worksNumber)) // no overflow

	for _, k := range ids[start:end] {
		idsString += fmt.Sprintf("&ids[]=%d", k)
	}

	return idsString, count, nil
}

func GetUserArtwork(c *fiber.Ctx, id, category string, page int) (User, error) {
	var user User

	token := session.GetPixivToken(c)

	URL := http.GetUserInformationURL(id)

	resp, err := http.UnwrapWebAPIRequest(c.Context(), URL, token)
	if err != nil {
		return user, err
	}

	resp = session.ProxyImageUrl(c, resp)

	err = json.Unmarshal([]byte(resp), &user)
	if err != nil {
		return user, err
	}

	if category != "bookmarks" {
		ids, count, err := GetUserArtworksID(c, id, category, page)
		if err != nil {
			return user, err
		}

		if count > 0 {
			// Check if the user has artworks available or not
			works, err := GetUserArtworks(c, id, ids)
			if err != nil {
				return user, err
			}

			// IDK but the order got shuffled even though Pixiv sorted the IDs in the response
			sort.Slice(works[:], func(i, j int) bool {
				left := works[i].ID
				right := works[j].ID
				return numberGreaterThan(left, right)
			})
			user.Artworks = works

			user.FrequentTags, err = GetFrequentTags(c, ids)
			if err != nil {
				return user, err
			}
		}

		// Artworks count
		user.ArtworksCount = count

	} else {
		// Bookmarks
		works, count, err := GetUserBookmarks(c, id, "show", page)
		if err != nil {
			return user, err
		}

		user.Artworks = works

		// Public bookmarks count
		user.ArtworksCount = count

	}

	err = user.ParseSocial()
	if err != nil {
		return User{}, err
	}

	if user.Background != nil {
		user.BackgroundImage = user.Background["url"].(string)
	}

	return user, nil
}

func GetUserBookmarks(c *fiber.Ctx, id, mode string, page int) ([]ArtworkBrief, int, error) {
	page--

	URL := http.GetUserBookmarksURL(id, mode, page)

	resp, err := http.UnwrapWebAPIRequest(c.Context(), URL, "")
	if err != nil {
		return nil, -1, err
	}
	resp = session.ProxyImageUrl(c, resp)

	var body struct {
		Artworks []json.RawMessage `json:"works"`
		Total    int               `json:"total"`
	}

	err = json.Unmarshal([]byte(resp), &body)
	if err != nil {
		return nil, -1, err
	}

	artworks := make([]ArtworkBrief, len(body.Artworks))

	for index, value := range body.Artworks {
		var artwork ArtworkBrief

		err = json.Unmarshal([]byte(value), &artwork)
		if err != nil {
			artworks[index] = ArtworkBrief{
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

func numberGreaterThan(l, r string) bool {
	if len(l) > len(r) {
		return true
	}
	if len(l) < len(r) {
		return false
	}
	return l > r
}
