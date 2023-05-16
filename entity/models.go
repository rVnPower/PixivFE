package entity

import (
	"html/template"
	"time"
)

type Illust struct {
	ID        int           `json:"id"`
	Title     string        `json:"title"`
	Caption   template.HTML `json:"caption"`
	Artist    UserBrief     `json:"user"`
	Date      time.Time     `json:"create_date"`
	Pages     int           `json:"page_count"`
	Views     int           `json:"total_view"`
	Bookmarks int           `json:"total_bookmarks"`
	Tags      []Tag         `json:"tags"`
	Images    []Image
}

type Spotlight struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Thumbnail string `json:"thumbnail"`
	URL       string `json:"article_url"`
	Date      string `json:"publish_date"`
}

type Tag struct {
	Name           string `json:"name"`
	TranslatedName string `json:"translated_name"`
}

type UserBrief struct {
	ID      int               `json:"id"`
	Name    string            `json:"name"`
	Account string            `json:"account"`
	Avatar  map[string]string `json:"profile_image_urls"`
}

type Image struct {
	Small    string `json:"square_medium"`
	Medium   string `json:"medium"`
	Large    string `json:"large"`
	Original string `json:"original"`
}
