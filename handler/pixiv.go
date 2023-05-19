package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"pixivfe/models"
	"regexp"
	"sort"
	"strconv"
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
	ArtworkRelatedURL     = "https://www.pixiv.net/ajax/illust/%s/recommend/init?limit=30"
	UserInformationURL    = "https://www.pixiv.net/ajax/user/%s?full=1"
	UserArtworksURL       = "https://www.pixiv.net/ajax/user/%s/profile/all"
	UserArtworksFullURL   = "https://www.pixiv.net/ajax/user/%s/profile/illusts?work_category=illustManga&is_first_page=0&lang=en%s"
)

func ProxyImage(url string) string {
	regex := regexp.MustCompile(`i\.pximg\.net`)
	proxy := "px2.rainchan.win"

	return regex.ReplaceAllString(url, proxy)
}

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

		image.Small = ProxyImage(imageRaw.Urls["thumb_mini"])
		image.Medium = ProxyImage(imageRaw.Urls["small"])
		image.Large = ProxyImage(imageRaw.Urls["regular"])
		image.Original = ProxyImage(imageRaw.Urls["original"])

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

func (p *PixivClient) GetUserArtworksID(id string) (*string, error) {
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

	i := 0
	for _, k := range ids {
		if i < 30 {
			idsString += fmt.Sprintf("&ids[]=%d", k)
		} else {
			break
		}
		i++
	}

	return &idsString, nil
}

func (p *PixivClient) GetRelatedArtworks(id string) ([]models.IllustShort, error) {
	url := fmt.Sprintf(ArtworkRelatedURL, id)

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

func (p *PixivClient) GetUserArtworks(id string) ([]models.IllustShort, error) {
	ids, err := p.GetUserArtworksID(id)
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

func (p *PixivClient) GetUserInformation(id string) (*models.User, error) {
	s, _ := p.TextRequest(fmt.Sprintf(UserInformationURL, id))

	var user *models.User
	var pr models.PixivResponse

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
	works, _ := p.GetUserArtworks(id)
	user.Artworks = works

	// Avatar
	user.Avatar = ProxyImage(user.Avatar)

	// Background image
	if body.Background != nil {
		user.BackgroundImage = ProxyImage(body.Background["url"])
	}

	return user, nil
}
