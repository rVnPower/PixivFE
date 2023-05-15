package entity

import (
	"html/template"
	"regexp"
	"time"
)

type Illust struct {
	ID            int                            `json:"id"`
	Title         string                         `json:"title"`
	Caption       template.HTML                  `json:"caption"`
	Artist        UserBrief                      `json:"user"`
	Date          time.Time                      `json:"create_date"`
	Pages         int                            `json:"page_count"`
	Views         int                            `json:"total_view"`
	Bookmarks     int                            `json:"total_bookmarks"`
	SingleImage   map[string]string              `json:"meta_single_page"`
	MultipleImage []map[string]map[string]string `json:"meta_pages"`
	Tags          []Tag                          `json:"tags"`
	Images        map[string]string              `json:"image_urls"`
	ImagesAlt     []Image
}

func ProxyImages(url string) string {
	regex := regexp.MustCompile(`i\.pximg\.net`)
	proxy := "px2.rainchan.win"

	return regex.ReplaceAllString(url, proxy)

}

func (c *Illust) ParseImages() {
	var images []Image
	if c.Pages > 1 {
		for _, imageMap := range c.MultipleImage {
			var image Image
			image.Small = imageMap["image_urls"]["square_medium"]
			image.Medium = imageMap["image_urls"]["medium"]
			image.Large = imageMap["image_urls"]["large"]
			image.Original = imageMap["image_urls"]["original"]

			images = append(images, image)
		}
	} else {
		var image Image
		image.Small = c.Images["square_medium"]
		image.Medium = c.Images["medium"]
		image.Large = c.Images["large"]
		image.Original = c.Images["original"]

		images = append(images, image)
	}

	c.ImagesAlt = images
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
