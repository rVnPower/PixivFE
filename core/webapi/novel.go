package core

import "time"

type Novel struct {
	BookmarkCount  int         `json:"bookmarkCount"`
	CommentCount   int         `json:"commentCount"`
	MarkerCount    int         `json:"markerCount"`
	CreateDate     time.Time   `json:"createDate"`
	UploadDate     time.Time   `json:"uploadDate"`
	Description    string      `json:"description"`
	ID             string      `json:"id"`
	Title          string      `json:"title"`
	LikeCount      int         `json:"likeCount"`
	PageCount      int         `json:"pageCount"`
	UserID         string      `json:"userId"`
	UserName       string      `json:"userName"`
	ViewCount      int         `json:"viewCount"`
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
			Tag       string `json:"tag"`
			Locked    bool   `json:"locked"`
			Deletable bool   `json:"deletable"`
			UserID    string `json:"userId"`
			UserName  string `json:"userName"`
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

func GetNovelByID(id string) {

}
