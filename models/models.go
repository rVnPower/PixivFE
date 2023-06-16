package models

import (
	"html/template"
	"time"

	"encoding/json"
)

type PaginationData struct {
	PreviousPage int
	CurrentPage  int
	NextPage     int
}

type PixivResponse struct {
	Error   bool
	Message string
	Body    json.RawMessage
}

type RankingResponse struct {
	Artworks []RankedArtwork `json:"contents"`
	Mode     string          `json:"mode"`
	Content  string          `json:"content"`
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

// Pixiv gives us 5 types of an image. I don't need the mini one tho.
// PS: Where tf is my 360x360 image, Pixiv?
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

type FrequentTag struct {
	Name           string `json:"tag"`
	TranslatedName string `json:"tag_translation"`
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
	XRestrict       xRestrict     `json:"xRestrict"`
	AiType          aiType        `json:"aiType"`
	User            UserShort
	RecentWorks     []IllustShort
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

type Comment struct {
	AuthorID   string `json:"userId"`
	AuthorName string `json:"userName"`
	Avatar     string `json:"img"`
	Context    string `json:"comment"`
	Stamp      string `json:"stampId"`
	Date       string `json:"commentDate"`
}

type User struct {
	ID              string        `json:"userId"`
	Name            string        `json:"name"`
	Avatar          string        `json:"imageBig"`
	BackgroundImage string        `json:"background"`
	Following       int           `json:"following"`
	MyPixiv         int           `json:"mypixivCount"`
	Comment         template.HTML `json:"commentHtml"`
	Artworks        []IllustShort `json:"artworks"`
	ArtworksCount   int
	FrequentTags    []FrequentTag
}

type UserShort struct {
	ID     string `json:"userId"`
	Name   string `json:"name"`
	Avatar string `json:"imageBig"`
}

type RankedArtwork struct {
	ID           int    `json:"illust_id"`
	Title        string `json:"title"`
	Rank         int    `json:"rank"`
	Pages        string `json:"illust_page_count"`
	Image        string `json:"url"`
	ArtistID     int    `json:"user_id"`
	ArtistName   string `json:"user_name"`
	ArtistAvatar string `json:"profile_img"`
}

type TagDetail struct {
	Name            string            `json:"tag"`
	AlternativeName string            `json:"word"`
	Metadata        map[string]string `json:"pixpedia"`
}

type PopularArtworks struct {
	Permanent []IllustShort `json:"permanent"`
	Recent    []IllustShort `json:"recent"`
}

type SearchArtworks struct {
	Artworks []IllustShort `json:"data"`
	Total    int           `json:"total"`
}

type SearchResult struct {
	Artworks    SearchArtworks
	Popular     PopularArtworks `json:"popular"`
	RelatedTags []string        `json:"relatedTags"`
}
