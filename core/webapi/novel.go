package core

import (
	"time"

	session "codeberg.org/vnpower/pixivfe/v2/core/config"
	http "codeberg.org/vnpower/pixivfe/v2/core/http"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

type Novel struct {
	Bookmarks      int         `json:"bookmarkCount"`
	CommentCount   int         `json:"commentCount"`
	MarkerCount    int         `json:"markerCount"`
	CreateDate     time.Time   `json:"createDate"`
	UploadDate     time.Time   `json:"uploadDate"`
	Description    string      `json:"description"`
	ID             string      `json:"id"`
	Title          string      `json:"title"`
	Likes          int         `json:"likeCount"`
	Pages          int         `json:"pageCount"`
	UserID         string      `json:"userId"`
	UserName       string      `json:"userName"`
	Views          int         `json:"viewCount"`
	IsOriginal     bool        `json:"isOriginal"`
	IsBungei       bool        `json:"isBungei"`
	XRestrict      int         `json:"xRestrict"`
	Restrict       int         `json:"restrict"`
	Content        string      `json:"content"`
	CoverURL       string      `json:"coverUrl"`
	IsBookmarkable bool        `json:"isBookmarkable"`
	BookmarkData   interface{} `json:"bookmarkData"`
	LikeData       bool        `json:"likeData"`
	PollData       interface{} `json:"pollData"`
	Marker         interface{} `json:"marker"`
	Tags           struct {
		AuthorID string `json:"authorId"`
		IsLocked bool   `json:"isLocked"`
		Tags     []struct {
			Name string `json:"tag"`
		} `json:"tags"`
		Writable bool `json:"writable"`
	} `json:"tags"`
	SeriesNavData  interface{} `json:"seriesNavData"`
	HasGlossary    bool        `json:"hasGlossary"`
	IsUnlisted     bool        `json:"isUnlisted"`
	Language       string      `json:"language"`
	CommentOff     int         `json:"commentOff"`
	CharacterCount int         `json:"characterCount"`
	WordCount      int         `json:"wordCount"`
	UseWordCount   bool        `json:"useWordCount"`
	ReadingTime    int         `json:"readingTime"`
	AiType         int         `json:"aiType"`
	Genre          string      `json:"genre"`
}

type NovelBrief struct {
	ID             string      `json:"id"`
	Title          string      `json:"title"`
	XRestrict      int         `json:"xRestrict"`
	Restrict       int         `json:"restrict"`
	CoverURL       string      `json:"url"`
	Tags           []string    `json:"tags"`
	UserID         string      `json:"userId"`
	UserName       string      `json:"userName"`
	UserAvatar     string      `json:"profileImageUrl"`
	TextCount      int         `json:"textCount"`
	WordCount      int         `json:"wordCount"`
	ReadingTime    int         `json:"readingTime"`
	Description    string      `json:"description"`
	IsBookmarkable bool        `json:"isBookmarkable"`
	BookmarkData   interface{} `json:"bookmarkData"`
	Bookmarks      int         `json:"bookmarkCount"`
	IsOriginal     bool        `json:"isOriginal"`
	CreateDate     time.Time   `json:"createDate"`
	UpdateDate     time.Time   `json:"updateDate"`
	IsMasked       bool        `json:"isMasked"`
	SeriesID       string      `json:"seriesId"`
	SeriesTitle    string      `json:"seriesTitle"`
	IsUnlisted     bool        `json:"isUnlisted"`
	AiType         int         `json:"aiType"`
	Genre          string      `json:"genre"`
}

func GetNovelByID(c *fiber.Ctx, id string) (Novel, error) {
	var novel Novel

	URL := http.GetNovelURL(id)

	response, err := http.UnwrapWebAPIRequest(c.Context(), URL, "")
	if err != nil {
		return novel, err
	}
	response = session.ProxyImageUrl(c, response)

	err = json.Unmarshal([]byte(response), &novel)
	if err != nil {
		return novel, err
	}

	return novel, nil
}
