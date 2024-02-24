package core

import (
	"strings"

	session "codeberg.org/vnpower/pixivfe/v2/core/user"
	http "codeberg.org/vnpower/pixivfe/v2/core/http"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

type TagDetail struct {
	Name            string `json:"tag"`
	AlternativeName string `json:"word"`
	Metadata        struct {
		Detail string      `json:"abstract"`
		Image  string      `json:"image"`
		Name   string      `json:"tag"`
		ID     json.Number `json:"id"`
	} `json:"pixpedia"`
}

type SearchArtworks struct {
	Artworks []ArtworkBrief `json:"data"`
	Total    int            `json:"total"`
}

type SearchResult struct {
	Artworks SearchArtworks
	Popular  struct {
		Permanent []ArtworkBrief `json:"permanent"`
		Recent    []ArtworkBrief `json:"recent"`
	} `json:"popular"`
	RelatedTags []string `json:"relatedTags"`
}

func GetTagData(c *fiber.Ctx, name string) (TagDetail, error) {
	var tag TagDetail

	URL := http.GetTagDetailURL(name)

	response, err := http.UnwrapWebAPIRequest(c.Context(), URL, "")
	if err != nil {
		return tag, err
	}

	response = session.ProxyImageUrl(c, response)

	err = json.Unmarshal([]byte(response), &tag)
	if err != nil {
		return tag, err
	}

	return tag, nil
}

func GetSearch(c *fiber.Ctx, artworkType, name, order, age_settings, page string) (*SearchResult, error) {

	URL := http.GetSearchArtworksURL(artworkType, name, order, age_settings, page)

	response, err := http.UnwrapWebAPIRequest(c.Context(), URL, "")
	if err != nil {
		return nil, err
	}
	response = session.ProxyImageUrl(c, response)

	// IDK how to do better than this lol
	temp := strings.ReplaceAll(string(response), `"illust"`, `"works"`)
	temp = strings.ReplaceAll(temp, `"manga"`, `"works"`)
	temp = strings.ReplaceAll(temp, `"illustManga"`, `"works"`)

	var resultRaw struct {
		*SearchResult
		ArtworksRaw json.RawMessage `json:"works"`
	}
	var artworks SearchArtworks
	var result *SearchResult

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
