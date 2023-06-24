package models

import (
	"regexp"
	"strings"
)

func ProxyImage(URL string, target string) string {

	if strings.Contains(URL, "s.pximg.net") {
		// This subdomain didn't get proxied
		return URL
	}

	regex := regexp.MustCompile(`.*?pximg\.net`)
	proxy := "https://" + target

	return regex.ReplaceAllString(URL, proxy)
}

func ProxyShortArtworkSlice(artworks []IllustShort, proxy string) []IllustShort {
	for i := range artworks {
		artworks[i].Thumbnail = ProxyImage(artworks[i].Thumbnail, proxy)
		artworks[i].ArtistAvatar = ProxyImage(artworks[i].ArtistAvatar, proxy)
	}

	return artworks
}

func ProxyCommentsSlice(comments []Comment, proxy string) []Comment {
	for i := range comments {
		comments[i].Avatar = ProxyImage(comments[i].Avatar, proxy)
	}
	return comments
}
