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
	Artworks    []RankedArtwork `json:"contents"`
	Mode        string          `json:"mode"`
	Content     string          `json:"content"`
	CurrentDate string          `json:"date"`
	PrevDateRaw json.RawMessage `json:"prev_date"`
	NextDateRaw json.RawMessage `json:"next_date"`
	PrevDate    string
	NextDate    string
}

func (s *RankingResponse) ProxyImages(proxy string) {
	s.Artworks = ProxyRankedArtworkSlice(s.Artworks, proxy)
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
	SanityLevel     int           `json:"sl"`
	XRestrict       xRestrict     `json:"xRestrict"`
	AiType          aiType        `json:"aiType"`
	User            UserShort
	RecentWorks     []IllustShort
	RelatedWorks    []IllustShort
	CommentsList    []Comment
	IsUgoira        bool
}

func (s *Illust) ProxyImages(proxy string) {
	for i := range s.Images {
		s.Images[i].Small = ProxyImage(s.Images[i].Small, proxy)
		s.Images[i].Medium = ProxyImage(s.Images[i].Medium, proxy)
		s.Images[i].Large = ProxyImage(s.Images[i].Large, proxy)
		s.Images[i].Original = ProxyImage(s.Images[i].Original, proxy)
	}
	for i := range s.RecentWorks {
		s.RecentWorks[i].Thumbnail = ProxyImage(s.RecentWorks[i].Thumbnail, proxy)
	}
	s.RelatedWorks = ProxyShortArtworkSlice(s.RelatedWorks, proxy)
	s.CommentsList = ProxyCommentsSlice(s.CommentsList, proxy)
	s.User.Avatar = ProxyImage(s.User.Avatar, proxy)
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

func (s *User) ProxyImages(proxy string) {
	s.Avatar = ProxyImage(s.Avatar, proxy)
	s.BackgroundImage = ProxyImage(s.BackgroundImage, proxy)
	s.Artworks = ProxyShortArtworkSlice(s.Artworks, proxy)
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
	Artworks []IllustShort `json:"data"`
	Total    int           `json:"total"`
}

type SearchResult struct {
	Artworks SearchArtworks
	Popular  struct {
		Permanent []IllustShort `json:"permanent"`
		Recent    []IllustShort `json:"recent"`
	} `json:"popular"`
	RelatedTags []string `json:"relatedTags"`
}

func (s *SearchResult) ProxyImages(proxy string) {
	s.Artworks.Artworks = ProxyShortArtworkSlice(s.Artworks.Artworks, proxy)
	s.Popular.Permanent = ProxyShortArtworkSlice(s.Popular.Permanent, proxy)
	s.Popular.Recent = ProxyShortArtworkSlice(s.Popular.Recent, proxy)
}

type Pixivision struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Thumbnail string `json:"thumbnailUrl"`
	URL       string `json:"url"`
}

type LandingRecommendByTags struct {
	Name     string `json:"tag"`
	Artworks []IllustShort
}

type LandingArtworks struct {
	Commissions     []IllustShort
	Following       []IllustShort
	Recommended     []IllustShort
	Newest          []IllustShort
	Rankings        []IllustShort
	Users           []IllustShort
	Pixivision      []Pixivision
	RecommendByTags []LandingRecommendByTags
}
