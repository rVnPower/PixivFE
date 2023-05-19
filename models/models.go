package models

import (
	"html/template"
	"time"

	"encoding/json"
)

type PixivResponse struct {
	Error   bool
	Message string
	Body    json.RawMessage
}

type ImageResponse struct {
	Urls map[string]string `json:"urls"`
}

type TagResponse struct {
	AuthorID string          `json:"authorId"`
	RawTags  json.RawMessage `json:"tags"`
}

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
	Safe: "Safe",
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
	NotAI:   "NotAI",
	AI:      "AI",
}

// Pixiv gives us 5 types of an image. I don't need the mini one tho.
// PS: Where tf is my 360x360 image, Pixiv? :rage:
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

type Illust struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Description template.HTML `json:"description"`
	UserID      string        `json:"userId"`
	UserName    string        `json:"userName"`
	UserAccount string        `json:"userAccount"`
	Date        time.Time     `json:"uploadDate"`
	Images      []Image       `json:"images"`
	Tags        []Tag         `json:"tags"`
	Pages       int           `json:"pageCount"`
	Bookmarks   int           `json:"bookmarkCount"`
	Likes       int           `json:"likeCount"`
	Comments    int           `json:"commentCount"`
	Views       int           `json:"viewCount"`
	XRestrict   xRestrict     `json:"xRestrict"`
	AiType      aiType        `json:"aiType"`
}

type IllustShort struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Description  template.HTML `json:"description"`
	ArtistID     string        `json:"userId"`
	ArtistName   string        `json:"userName"`
	ArtistAvatar string        `json:"profileImageUrl"`
	Date         time.Time     `json:"uploadDate"`
	Thumbnail    string        `json:"url"`
	Pages        int           `json:"pageCount"`
	XRestrict    xRestrict     `json:"xRestrict"`
	AiType       aiType        `json:"aiType"`
}

type User struct {
	ID              string        `json:"userId"`
	Name            string        `json:"name"`
	Avatar          string        `json:"imageBig"`
	BackgroundImage string        `json:"background"`
	Following       int           `json:"following"`
	MyPixiv         int           `json:"mypixivCount"`
	Comment         string        `json:"comment"`
	Artworks        []IllustShort `json:"artworks"`
}
