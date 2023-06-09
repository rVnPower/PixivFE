package utils

import (
	"pixivfe/configs"
	"regexp"
	"strings"
)

func ProxyImage(url string) string {
	if strings.Contains(url, "s.pximg.net") {
		// This subdomain didn't get proxied
		return url
	}

	regex := regexp.MustCompile(`.*?pximg\.net`)
	proxy := "https://" + configs.ProxyServer

	return regex.ReplaceAllString(url, proxy)
}
