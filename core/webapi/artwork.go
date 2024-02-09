package core

import (
	"fmt"
	"html/template"
	"sort"
	"strconv"
	"strings"
	"time"

	session "codeberg.org/vnpower/pixivfe/v2/core/config"
	http "codeberg.org/vnpower/pixivfe/v2/core/http"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

// Pixiv returns 0, 1, 2 to filter SFW and/or NSFW artworks.
// Those values are saved in `xRestrict`
// 0: Safe
// 1: R18
// 2: R18G
type xRestrict int

const (
	Safe xRestrict = 0
	R18  xRestrict = 1
	R18G xRestrict = 2
)

var xRestrictModel = map[xRestrict]string{
	Safe: "",
	R18:  "R18",
	R18G: "R18G",
}

// Pixiv returns 0, 1, 2 to filter SFW and/or NSFW artworks.
// Those values are saved in `aiType`
// 0: Not rated / Unknown
// 1: Not AI-generated
// 2: AI-generated

type aiType int

const (
	Unrated aiType = 0
	NotAI   aiType = 1
	AI      aiType = 2
)

var aiTypeModel = map[aiType]string{
	Unrated: "Unrated",
	NotAI:   "Not AI",
	AI:      "AI",
}

type ImageResponse struct {
	Urls map[string]string `json:"urls"`
}

type Image struct {
	Small    string `json:"thumb_mini"`
	Medium   string `json:"small"`
	Large    string `json:"regular"`
	Original string `json:"original"`
}

type Tag struct {
	Name           string `json:"tag"`
	TranslatedName string `json:"translation"`
}

type Comment struct {
	AuthorID   string `json:"userId"`
	AuthorName string `json:"userName"`
	Avatar     string `json:"img"`
	Context    string `json:"comment"`
	Stamp      string `json:"stampId"`
	Date       string `json:"commentDate"`
}

type UserBrief struct {
	ID     string `json:"userId"`
	Name   string `json:"name"`
	Avatar string `json:"imageBig"`
}

type Illust struct {
	ID              string        `json:"id"`
	Title           string        `json:"title"`
	Description     template.HTML `json:"description"`
	UserID          string        `json:"userId"`
	UserName        string        `json:"userName"`
	UserAccount     string        `json:"userAccount"`
	Date            time.Time     `json:"uploadDate"`
	Images          []Image       `json:"images"`
	Tags            []Tag         `json:"tags"`
	Pages           int           `json:"pageCount"`
	Bookmarks       int           `json:"bookmarkCount"`
	Likes           int           `json:"likeCount"`
	Comments        int           `json:"commentCount"`
	Views           int           `json:"viewCount"`
	CommentDisabled int           `json:"commentOff"`
	SanityLevel     int           `json:"sl"`
	XRestrict       xRestrict     `json:"xRestrict"`
	AiType          aiType        `json:"aiType"`
	Bookmarked      any           `json:"bookmarkData"`
	Liked           any           `json:"json:"likeData"`
	User            UserBrief
	RecentWorks     []ArtworkBrief
	RelatedWorks    []ArtworkBrief
	CommentsList    []Comment
	IsUgoira        bool
}

func GetUserBasicInformation(c *fiber.Ctx, id string) (UserBrief, error) {
	var user UserBrief

	URL := http.GetUserInformationURL(id)

	response, err := http.UnwrapWebAPIRequest(URL, "")
	if err != nil {
		return user, err
	}
	response = session.ProxyImageUrl(response)

	err = json.Unmarshal([]byte(response), &user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func GetArtworkImages(c *fiber.Ctx, id string) ([]Image, error) {
	var resp []ImageResponse
	var images []Image

	URL := http.GetArtworkImagesURL(id)

	response, err := http.UnwrapWebAPIRequest(URL, "")
	if err != nil {
		return nil, err
	}
	response = session.ProxyImageUrl(response)

	err = json.Unmarshal([]byte(response), &resp)
	if err != nil {
		return images, err
	}

	// Extract and proxy every images
	for _, imageRaw := range resp {
		var image Image

		image.Small = imageRaw.Urls["thumb_mini"]
		image.Medium = imageRaw.Urls["small"]
		image.Large = imageRaw.Urls["regular"]
		image.Original = imageRaw.Urls["original"]

		images = append(images, image)
	}

	return images, nil
}

func GetArtworkComments(c *fiber.Ctx, id string) ([]Comment, error) {
	var body struct {
		Comments []Comment `json:"comments"`
	}

	URL := http.GetArtworkCommentsURL(id)

	response, err := http.UnwrapWebAPIRequest(URL, "")
	if err != nil {
		return nil, err
	}
	response = session.ProxyImageUrl(response)

	err = json.Unmarshal([]byte(response), &body)
	if err != nil {
		return nil, err
	}

	return body.Comments, nil
}

func GetRelatedArtworks(c *fiber.Ctx, id string) ([]ArtworkBrief, error) {
	var body struct {
		Illusts []ArtworkBrief `json:"illusts"`
	}

	// TODO: keep the hard-coded limit?
	URL := http.GetArtworkRelatedURL(id, 96)

	response, err := http.UnwrapWebAPIRequest(URL, "")
	if err != nil {
		return nil, err
	}

	response = session.ProxyImageUrl(response)

	err = json.Unmarshal([]byte(response), &body)
	if err != nil {
		return nil, err
	}

	return body.Illusts, nil
}

func GetArtworkByID(c *fiber.Ctx, id string) (*Illust, error) {
	var images []Image

	URL := http.GetArtworkInformationURL(id)

	response, err := http.UnwrapWebAPIRequest(URL, "")
	if err != nil {
		return nil, err
	}

	var illust struct {
		*Illust

		Recent  map[int]any     `json:"userIllusts"`
		RawTags json.RawMessage `json:"tags"`
	}

	// Parse basic illust information
	err = json.Unmarshal([]byte(response), &illust)
	if err != nil {
		return nil, err
	}

	// Begin testing here

	c1 := make(chan []Image)
	c2 := make(chan []ArtworkBrief)
	c3 := make(chan UserBrief)
	c4 := make(chan []Tag)
	c5 := make(chan []ArtworkBrief)
	c6 := make(chan []Comment)

	go func() {
		// Get illust images
		images, err = GetArtworkImages(c, id)
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

		recent, err := GetUserArtworks(c, illust.UserID, idsString)
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
		userInfo, err := GetUserBasicInformation(c, illust.UserID)
		if err != nil {
			//
		}
		c3 <- userInfo
	}()

	go func() {
		var tagsList []Tag
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
			var newTag Tag
			newTag.Name = tag.Tag
			newTag.TranslatedName = tag.Translation["en"]

			tagsList = append(tagsList, newTag)
		}
		c4 <- tagsList
	}()

	go func() {
		related, _ := GetRelatedArtworks(c, id)
		// Error handling...
		c5 <- related
	}()

	go func() {
		comments, _ := GetArtworkComments(c, id)
		// Error handling...
		c6 <- comments
	}()

	illust.Images = <-c1
	illust.RecentWorks = <-c2
	illust.User = <-c3
	illust.Tags = <-c4
	illust.RelatedWorks = <-c5
	illust.CommentsList = <-c6

	// If this artwork is an ugoira
	illust.IsUgoira = strings.Contains(illust.Images[0].Original, "ugoira")

	return illust.Illust, nil
}
