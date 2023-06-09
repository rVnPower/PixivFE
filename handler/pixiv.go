package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"math"
	"net/http"
	"pixivfe/models"
	"pixivfe/utils"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
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
	ArtworkDiscoveryURL   = "https://www.pixiv.net/ajax/discovery/artworks?mode=%s&limit=%d"
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

	if resp.StatusCode != 200 {
		return resp, errors.New(fmt.Sprintf("Pixiv returned code: %d for request ", resp.StatusCode))
	}

	return resp, nil
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func ParseJSONToString(json, key string) string {
	return gjson.Get(json, key).String()
}

func ParseJSONToInt(json, key string) int {
	return int(gjson.Get(json, key).Int())
}

func ParseJSON(json, key string) gjson.Result {
	return gjson.Get(json, key)
}

func ExtractImageURL(data string, mode int) models.Image {
	// 1 => thumb_mini, small, regular, original
	// 2 => mini, thumb, small, regular, original
	var image models.Image

	switch mode {
	case 1:
		image.Small = utils.ProxyImage(ParseJSONToString(data, "thumb_mini"))
		image.Medium = utils.ProxyImage(ParseJSONToString(data, "small"))
		image.Large = utils.ProxyImage(ParseJSONToString(data, "regular"))
		image.Original = utils.ProxyImage(ParseJSONToString(data, "original"))
	case 2:
		image.Small = utils.ProxyImage(ParseJSONToString(data, "thumb"))
		image.Medium = utils.ProxyImage(ParseJSONToString(data, "small"))
		image.Large = utils.ProxyImage(ParseJSONToString(data, "regular"))
		image.Original = utils.ProxyImage(ParseJSONToString(data, "original"))
	}

	return image
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

	var images []models.Image

	if ParseJSON(s, "error").Bool() {
		return images, errors.New(fmt.Sprintf("Pixiv returned error message: %s", ParseJSONToString(s, "message")))
	}

	// Extract and proxy every images
	for _, imageRaw := range gjson.Get(s, "body.#.urls").Array() {
		data := imageRaw.String()
		image := ExtractImageURL(data, 1)

		images = append(images, image)
	}

	return images, nil
}

func (p *PixivClient) GetArtworkByID(id string) (models.Illust, error) {
	s, _ := p.TextRequest(fmt.Sprintf(ArtworkInformationURL, id))

	var artwork models.Illust

	if ParseJSON(s, "error").Bool() {
		return artwork, errors.New(fmt.Sprintf("Pixiv returned error message: %s", ParseJSONToString(s, "message")))
	}

	body := ParseJSONToString(s, "body")

	artwork.ID = ParseJSONToString(body, "id")
	artwork.Title = ParseJSONToString(body, "title")
	artwork.Description = template.HTML(ParseJSONToString(body, "description"))
	artwork.UserID = ParseJSONToString(body, "userId")
	artwork.UserName = ParseJSONToString(body, "userName")
	artwork.UserAccount = ParseJSONToString(body, "userAccount")
	artwork.Date, _ = time.Parse("2006-01-02T03:04:00+00:00", ParseJSONToString(body, "uploadDate"))
	artwork.Images, _ = p.GetArtworkImages(id)
	artwork.Pages = ParseJSONToInt(body, "pageCount")
	artwork.Bookmarks = ParseJSONToInt(body, "bookmarkCount")
	artwork.Likes = ParseJSONToInt(body, "likeCount")
	artwork.Comments = ParseJSONToInt(body, "commentCount")
	artwork.Views = ParseJSONToInt(body, "viewCount")
	artwork.XRestrict = ParseJSONToInt(body, "xRestrict")
	artwork.AiType = ParseJSONToInt(body, "aiType")
	artwork.Tags = make([]models.Tag, 0)

	for _, rawTag := range ParseJSON(body, "tags.tags").Array() {
		var tag models.Tag
		data := rawTag.String()

		tag.Name = ParseJSONToString(data, "tag")
		tag.TranslatedName = ParseJSONToString(data, "tag.translation.en")

		artwork.Tags = append(artwork.Tags, tag)
	}

	return artwork, nil
}

func (p *PixivClient) GetArtworkComments(id string) ([]models.Comment, error) {
	comments := make([]models.Comment, 0)

	s, _ := p.TextRequest(fmt.Sprintf(ArtworkCommentsURL, id))

	if ParseJSON(s, "error").Bool() {
		return comments, errors.New(fmt.Sprintf("Pixiv returned error message: %s", ParseJSONToString(s, "message")))
	}

	for _, rawComment := range ParseJSON(s, "body.comments").Array() {
		data := rawComment.String()
		var comment models.Comment

		comment.AuthorID = ParseJSONToString(data, "userId")
		comment.AuthorName = ParseJSONToString(data, "userName")
		comment.Avatar = utils.ProxyImage(ParseJSONToString(data, "img"))
		comment.Context = ParseJSONToString(data, "comment")
		comment.Stamp = ParseJSONToString(data, "stampId")
		comment.Date = ParseJSONToString(data, "commentDate")

		comments = append(comments, comment)
	}

	return comments, nil
}

func (p *PixivClient) GetUserArtworksID(id string, page int) (*string, error) {
	s, _ := p.TextRequest(fmt.Sprintf(UserArtworksURL, id))

	var pr models.PixivResponse

	err := json.Unmarshal([]byte(s), &pr)

	if err != nil {
		return nil, err
	}
	if pr.Error {
		return nil, errors.New(fmt.Sprintf("Pixiv returned error message: %s", pr.Message))
	}

	var ids []int
	var idsString string
	var body struct {
		Illusts map[int]string `json:"illusts"`
	}
	err = json.Unmarshal(pr.Body, &body)
	if err != nil {
		return nil, err
	}

	// Get the keys, because Pixiv only returns IDs (very evil)
	for k := range body.Illusts {
		ids = append(ids, k)
	}

	// Reverse sort the ids
	sort.Sort(sort.Reverse(sort.IntSlice(ids)))

	worksNumber := float64(len(ids))
	worksPerPage := 30.0

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
		return -1, err
	}
	if pr.Error {
		return -1, errors.New(fmt.Sprintf("Pixiv returned error message: %s", pr.Message))
	}

	var body struct {
		Illusts map[int]string `json:"illusts"`
	}
	err = json.Unmarshal(pr.Body, &body)
	if err != nil {
		return -1, err
	}

	return len(body.Illusts), nil
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

func (p *PixivClient) GetUserArtworks(id string, page int) ([]models.IllustShort, error) {
	ids, err := p.GetUserArtworksID(id, page)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(UserArtworksFullURL, id, *ids)

	var pr models.PixivResponse
	var works []models.IllustShort

	s, _ := p.TextRequest(url)

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
	works, _ := p.GetUserArtworks(id, page)
	user.Artworks = works

	// Background image
	if body.Background != nil {
		user.BackgroundImage = body.Background["url"].(string)
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
		if err != nil {
			return nil, err
		}

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

	return newWorks, nil
}

func (p *PixivClient) GetRanking(mode string, content string, page string) (models.RankingResponse, error) {
	// Ranking data is formatted differently
	var pr models.RankingResponse

	url := fmt.Sprintf(ArtworkRankingURL, mode, content, page)

	s, err := p.TextRequest(url)

	if err != nil {
		return pr, err
	}

	err = json.Unmarshal([]byte(s), &pr)
	if err != nil {
		return pr, err
	}

	return pr, nil
}

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

func (p *PixivClient) GetDiscoveryArtwork(mode string, count int) ([]models.IllustShort, error) {
	var artworks []models.IllustShort

	for count > 0 {
		var pr models.PixivResponse
		itemsForRequest := Min(100, count)

		count -= itemsForRequest

		url := fmt.Sprintf(ArtworkDiscoveryURL, mode, itemsForRequest)
		s, err := p.TextRequest(url)

		if err != nil {
			return artworks, err
		}

		err = json.Unmarshal([]byte(s), &pr)

		if pr.Error {
			return artworks, errors.New(pr.Message)
		}

		var thumbnail struct {
			Data json.RawMessage `json:"thumbnails"`
		}

		err = json.Unmarshal([]byte(pr.Body), &thumbnail)
		if err != nil {
			return nil, err
		}

		var body struct {
			Artworks []models.IllustShort `json:"illust"`
		}

		err = json.Unmarshal([]byte(thumbnail.Data), &body)
		if err != nil {
			return nil, err
		}

		artworks = append(artworks, body.Artworks...)
	}

	return artworks, nil
}
