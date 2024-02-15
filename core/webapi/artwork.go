package core

import (
	"errors"
	"fmt"
	"html/template"
	"sort"
	"strconv"
	"strings"
	"sync"
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
	Width  int               `json:"width"`
	Height int               `json:"height"`
	Urls   map[string]string `json:"urls"`
}

type Image struct {
	Width    int
	Height   int
	Small    string
	Medium   string
	Large    string
	Original string
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

type ArtworkBrief struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	ArtistID     string `json:"userId"`
	ArtistName   string `json:"userName"`
	ArtistAvatar string `json:"profileImageUrl"`
	Thumbnail    string `json:"url"`
	Pages        int    `json:"pageCount"`
	XRestrict    int    `json:"xRestrict"`
	AiType       int    `json:"aiType"`
	Bookmarked   any    `json:"bookmarkData"`
	IllustType   int    `json:"illustType"`
}

type Illust struct {
	ID              string        `json:"id"`
	Title           string        `json:"title"`
	Description     template.HTML `json:"description"`
	UserID          string        `json:"userId"`
	UserName        string        `json:"userName"`
	UserAccount     string        `json:"userAccount"`
	Date            time.Time     `json:"uploadDate"`
	Images          []Image
	Tags            []Tag     `json:"tags"`
	Pages           int       `json:"pageCount"`
	Bookmarks       int       `json:"bookmarkCount"`
	Likes           int       `json:"likeCount"`
	Comments        int       `json:"commentCount"`
	Views           int       `json:"viewCount"`
	CommentDisabled int       `json:"commentOff"`
	SanityLevel     int       `json:"sl"`
	XRestrict       xRestrict `json:"xRestrict"`
	AiType          aiType    `json:"aiType"`
	Bookmarked      any       `json:"bookmarkData"`
	Liked           any       `json:"likeData"`
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
	response = session.ProxyImageUrl(c, response)

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
	response = session.ProxyImageUrl(c, response)

	err = json.Unmarshal([]byte(response), &resp)
	if err != nil {
		return images, err
	}

	// Extract and proxy every images
	for _, imageRaw := range resp {
		var image Image

		// this is the original art dimention, not the "regular" art dimension
		// the image ratio of "regular" is close to Width/Height
		// maybe not useful
		image.Width = imageRaw.Width
		image.Height = imageRaw.Height

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
	response = session.ProxyImageUrl(c, response)

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

	response = session.ProxyImageUrl(c, response)

	err = json.Unmarshal([]byte(response), &body)
	if err != nil {
		return nil, err
	}

	return body.Illusts, nil
}

func GetArtworkByID(c *fiber.Ctx, id string, full bool) (*Illust, error) {
	URL := http.GetArtworkInformationURL(id)

	response, err := http.UnwrapWebAPIRequest(URL, "")
	if err != nil {
		return nil, err
	}

	var illust struct {
		*Illust

		// recent illustrations by same user
		Recent  map[int]any     `json:"userIllusts"`
		RawTags json.RawMessage `json:"tags"`
	}

	// Parse basic illust information
	err = json.Unmarshal([]byte(response), &illust)
	if err != nil {
		return nil, err
	}

	// Begin testing here

	wg := sync.WaitGroup{}
	cerr := make(chan error, 6)

	wg.Add(3)

	go func() {
		// Get illust images
		defer wg.Done()
		images, err := GetArtworkImages(c, id)
		if err != nil {

			cerr <- err
			return
		}
		illust.Images = images
	}()

	go func() {
		// Get basic user information (the URL above does not contain avatars)
		defer wg.Done()
		var err error
		userInfo, err := GetUserBasicInformation(c, illust.UserID)
		if err != nil {
			cerr <- err
			return
		}
		illust.User = userInfo
	}()

	go func() {
		defer wg.Done()
		var err error
		// Extract tags
		var tags struct {
			Tags []struct {
				Tag         string            `json:"tag"`
				Translation map[string]string `json:"translation"`
			} `json:"tags"`
		}
		err = json.Unmarshal(illust.RawTags, &tags)
		if err != nil {
			cerr <- err
			return
		}

		var tagsList []Tag
		for _, tag := range tags.Tags {
			var newTag Tag
			newTag.Name = tag.Tag
			newTag.TranslatedName = tag.Translation["en"]

			tagsList = append(tagsList, newTag)
		}
		illust.Tags = tagsList
	}()

	if full {
		wg.Add(3)

		go func() {
			defer wg.Done()
			var err error
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
				cerr <- err
				return
			}
			sort.Slice(recent[:], func(i, j int) bool {
				left, _ := strconv.Atoi(recent[i].ID)
				right, _ := strconv.Atoi(recent[j].ID)
				return left > right
			})
			illust.RecentWorks = recent
		}()

		go func() {
			defer wg.Done()
			var err error
			related, err := GetRelatedArtworks(c, id)
			if err != nil {
				cerr <- err
				return
			}
			illust.RelatedWorks = related
		}()

		go func() {
			defer wg.Done()
			if illust.CommentDisabled == 1 {
				return
			}
			var err error
			comments, err := GetArtworkComments(c, id)
			if err != nil {
				println("here")
				cerr <- err
				return
			}
			illust.CommentsList = comments
		}()
	}

	wg.Wait()
	close(cerr)

	all_errors := []error{}
	for suberr := range cerr {
		all_errors = append(all_errors, suberr)
	}
	err_summary := errors.Join(all_errors...)
	if err_summary != nil {
		return nil, err_summary
	}

	// If this artwork is an ugoira
	illust.IsUgoira = strings.Contains(illust.Images[0].Original, "ugoira")

	return illust.Illust, nil
}
