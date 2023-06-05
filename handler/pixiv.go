package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"pixivfe/configs"
	"pixivfe/models"
	"sort"
	"strconv"
	"strings"
)

type PixivClient struct {
	Client *http.Client

	Cookie map[string]string
	Header map[string]string
	Lang   string
}

const (
	ArtworkInformationURL = "https://www.pixiv.net/ajax/illust/%s"
	ArtworkImagesURL      = "https://www.pixiv.net/ajax/illust/%s/pages"
	ArtworkRelatedURL     = "https://www.pixiv.net/ajax/illust/%s/recommend/init?limit=%d"
	ArtworkCommentsURL    = "https://www.pixiv.net/ajax/illusts/comments/roots?illust_id=%s&limit=100"
	ArtworkNewestURL      = "https://www.pixiv.net/ajax/illust/new?limit=30&type=%s&r18=%s&lastId=%s"
	ArtworkRankingURL     = "https://www.pixiv.net/ranking.php?format=json&mode=%s&content=%s&p=%s"
	ArtworkDiscoveryURL   = "https://www.pixiv.net/ajax/discovery/artworks?mode=%s&limit=100"
	SearchTagURL          = "https://www.pixiv.net/ajax/search/tags/%s"
	SearchArtworksURL     = "https://www.pixiv.net/ajax/search/%s/%s?order=%s&mode=%s&p=%s"
	SearchTopURL          = "https://www.pixiv.net/ajax/search/top/%s"
	UserInformationURL    = "https://www.pixiv.net/ajax/user/%s?full=1"
	UserArtworksURL       = "https://www.pixiv.net/ajax/user/%s/profile/all"
	UserArtworksFullURL   = "https://www.pixiv.net/ajax/user/%s/profile/illusts?work_category=illustManga&is_first_page=0&lang=en%s"
)

func (p *PixivClient) SetHeader(header map[string]string) {
	p.Header = header
}

func (p *PixivClient) AddHeader(key, value string) {
	p.Header[key] = value
}

func (p *PixivClient) SetUserAgent(value string) {
	p.AddHeader("User-Agent", value)
}

func (p *PixivClient) SetCookie(cookie map[string]string) {
	p.Cookie = cookie
}

func (p *PixivClient) AddCookie(key, value string) {
	p.Cookie[key] = value
}

func (p *PixivClient) SetSessionID(value string) {
	p.Cookie["PHPSESSID"] = value
}

func (p *PixivClient) SetLang(lang string) {
	p.Lang = lang
}

func (p *PixivClient) Request(URL string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", URL, nil)

	// Add headers
	for k, v := range p.Header {
		req.Header.Add(k, v)
	}
	for k, v := range p.Cookie {
		req.AddCookie(&http.Cookie{Name: k, Value: v})
	}
	// Make a request
	resp, err := p.Client.Do(req)

	if err != nil {
		return resp, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return resp, errors.New("404 returned")
	}

	if resp.StatusCode != 200 {
		return resp, errors.New(fmt.Sprintf("Server returned code: %d", resp.StatusCode))
	}

	return resp, nil
}

func (p *PixivClient) TextRequest(URL string) (string, error) {
	resp, err := p.Request(URL)
	if err != nil {
		return "", err
	}

	// Extract the bytes from server's response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return string(body), err
	}

	return string(body), nil
}

func (p *PixivClient) GetArtworkImages(id string) ([]models.Image, error) {
	s, _ := p.TextRequest(fmt.Sprintf(ArtworkImagesURL, id))

	var pr models.PixivResponse
	var resp []models.ImageResponse
	var images []models.Image

	err := json.Unmarshal([]byte(s), &pr)
	if err != nil {
		return images, errors.New(fmt.Sprintf("Failed to extract data from JSON response from server. %s", err))
	}
	if pr.Error {
		return images, errors.New(fmt.Sprintf("Server refuses to return data. %s", err))
	}

	err = json.Unmarshal([]byte(pr.Body), &resp)
	if err != nil {
		return images, errors.New(fmt.Sprintf("Failed to extract images for illust. %s", err))
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

	var resp models.PixivResponse
	var images []models.Image

	// Parse Pixiv response body
	err := json.Unmarshal([]byte(s), &resp)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to extract data from JSON response from server. %s", err))
	}
	if resp.Error {
		return nil, errors.New(fmt.Sprintf("Server refuses to return data. %s", err))
	}

	var illust struct {
		*models.Illust
		RawTags json.RawMessage `json:"tags"`
	}

	// Parse basic illust information
	err = json.Unmarshal([]byte(resp.Body), &illust)

	// Get illust images
	images, err = p.GetArtworkImages(id)
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	illust.Images = images

	// Extract tags
	var tags struct {
		Tags []struct {
			Tag         string            `json:"tag"`
			Translation map[string]string `json:"translation"`
		} `json:"tags"`
	}
	err = json.Unmarshal(illust.RawTags, &tags)
	if err != nil {
		fmt.Printf("%s\n", err)
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
		return nil, errors.New(fmt.Sprintf("Failed to extract data from JSON response from server. %s", err))
	}
	if pr.Error {
		return nil, errors.New(pr.Message)
	}

	err = json.Unmarshal([]byte(pr.Body), &body)

	return body.Comments, nil
}

func (p *PixivClient) GetUserArtworksID(id string, page int) (*string, error) {
	s, _ := p.TextRequest(fmt.Sprintf(UserArtworksURL, id))

	var pr models.PixivResponse

	err := json.Unmarshal([]byte(s), &pr)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to extract data from JSON response from server. %s", err))
	}
	if pr.Error {
		return nil, errors.New(pr.Message)
	}

	var ids []int
	var idsString string
	var body struct {
		Illusts map[int]string `json:"illusts"`
	}
	err = json.Unmarshal(pr.Body, &body)

	// Get the keys, because Pixiv only returns IDs (very evil)
	for k := range body.Illusts {
		ids = append(ids, k)
	}

	// Reverse sort the ids
	sort.Sort(sort.Reverse(sort.IntSlice(ids)))

	worksNumber := float64(len(ids))
	worksPerPage := float64(configs.Configs.PageItems)

	if page < 1 || float64(page) > math.Ceil(worksNumber/worksPerPage)+1.0 {
		return nil, errors.New("Page overflow")
	}

	start := (page - 1) * int(worksPerPage)
	end := int(math.Min(float64(page)*worksPerPage, worksNumber)) // no overflow

	for _, k := range ids[start:end] {
		idsString += fmt.Sprintf("&ids[]=%d", k)
	}

	return &idsString, nil
}

func (p *PixivClient) GetUserArtworksCount(id string) (int, error) {
	s, _ := p.TextRequest(fmt.Sprintf(UserArtworksURL, id))

	var pr models.PixivResponse

	err := json.Unmarshal([]byte(s), &pr)

	if err != nil {
		return -1, errors.New(fmt.Sprintf("Failed to extract data from JSON response from server. %s", err))
	}
	if pr.Error {
		return -1, errors.New(pr.Message)
	}

	var body struct {
		Illusts map[int]string `json:"illusts"`
	}
	err = json.Unmarshal(pr.Body, &body)

	return len(body.Illusts), nil
}

func (p *PixivClient) GetRelatedArtworks(id string) ([]models.IllustShort, error) {
	url := fmt.Sprintf(ArtworkRelatedURL, id, configs.Configs.PageItems)

	var pr models.PixivResponse

	s, err := p.TextRequest(url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(s), &pr)

	var body struct {
		Illusts []models.IllustShort `json:"illusts"`
	}

	err = json.Unmarshal([]byte(pr.Body), &body)
	if err != nil {
		return nil, err
	}

	return body.Illusts, nil
}

func (p *PixivClient) GetUserArtworks(id string, page int) ([]models.IllustShort, error) {
	ids, err := p.GetUserArtworksID(id, page)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf(UserArtworksFullURL, id, *ids)

	var pr models.PixivResponse
	var works []models.IllustShort

	s, err := p.TextRequest(url)

	err = json.Unmarshal([]byte(s), &pr)

	var body struct {
		Illusts map[int]json.RawMessage `json:"works"`
	}

	err = json.Unmarshal(pr.Body, &body)

	for _, v := range body.Illusts {
		var illust models.IllustShort
		err = json.Unmarshal(v, &illust)

		works = append(works, illust)
	}

	// IDK but the order got shuffled even though Pixiv sorted the IDs in the response
	sort.Slice(works[:], func(i, j int) bool {
		left, _ := strconv.Atoi(works[i].ID)
		right, _ := strconv.Atoi(works[j].ID)
		return left > right
	})

	return works, nil
}

func (p *PixivClient) GetUserInformation(id string, page int) (*models.User, error) {
	var user *models.User
	var pr models.PixivResponse

	s, _ := p.TextRequest(fmt.Sprintf(UserInformationURL, id))

	err := json.Unmarshal([]byte(s), &pr)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to extract data from JSON response from server. %s", err))
	}
	if pr.Error {
		return nil, errors.New(pr.Message)
	}

	var body struct {
		*models.User
		Background map[string]string `json:"background"`
	}

	// Basic user information
	err = json.Unmarshal([]byte(pr.Body), &body)

	user = body.User

	// Artworks
	works, _ := p.GetUserArtworks(id, page)
	user.Artworks = works

	// Background image
	if body.Background != nil {
		user.BackgroundImage = body.Background["url"]
	}

	// Artworks count
	// user.ArtworksCount, _ = p.GetUserArtworksCount(id)

	return user, nil
}

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

	println(len(newWorks))

	return newWorks, nil
}

func (p *PixivClient) GetRanking(mode string, content string, page string) (models.RankingResponse, error) {
	// Ranking data is formatted differently
	var pr models.RankingResponse

	url := fmt.Sprintf(ArtworkRankingURL, mode, content, page)

	s, err := p.TextRequest(url)

	err = json.Unmarshal([]byte(s), &pr)

	_ = err

	return pr, nil
}

func (p *PixivClient) GetTagData(name string) (models.TagDetail, error) {
	var pr models.PixivResponse
	var tag models.TagDetail

	url := fmt.Sprintf(SearchTagURL, name)

	s, err := p.TextRequest(url)

	err = json.Unmarshal([]byte(s), &pr)

	if err != nil {
		return tag, errors.New("Error")
	}

	if pr.Error {
		return tag, errors.New(pr.Message)
	}

	err = json.Unmarshal([]byte(pr.Body), &tag)

	return tag, nil
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

	result = resultRaw.SearchResult

	err = json.Unmarshal([]byte(resultRaw.ArtworksRaw), &artworks)

	result.Artworks = artworks

	return result, nil
}

func (p *PixivClient) GetDiscoveryArtwork(mode string) ([]models.IllustShort, error) {
	var pr models.PixivResponse
	var artworks []models.IllustShort

	url := fmt.Sprintf(ArtworkDiscoveryURL, mode)

	s, err := p.TextRequest(url)

	err = json.Unmarshal([]byte(s), &pr)

	if err != nil {
		return artworks, errors.New("Error")
	}

	if pr.Error {
		return artworks, errors.New(pr.Message)
	}

	var thumbnail struct {
		Data json.RawMessage `json:"thumbnails"`
	}

	err = json.Unmarshal([]byte(pr.Body), &thumbnail)

	var body struct {
		Artworks []models.IllustShort `json:"illust"`
	}

	err = json.Unmarshal([]byte(thumbnail.Data), &body)

	return body.Artworks, nil
}
