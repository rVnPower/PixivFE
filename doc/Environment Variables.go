// Environment Variables
//
// PixivFE's behavior is governed by those Environment Variables.

package doc

import (
	"log"
	"os"
)

// An environment variable is a KEY=VALUE pair
type EnvVar = struct {
	Name       string
	CommonName string
	Value      string // available at run-time
	Announce   bool
}

// All environment variables used by PixivFE
var EnvList []*EnvVar = []*EnvVar{
	{
		Name:       "PIXIVFE_DEV",
		CommonName: "development mode",
		// **Required**: No
		//
		// Set this to anything to enable development mode, in which the server will live-reload HTML templates and disable caching. For example, `PIXIVFE_DEV=true`.
	},

	{
		Name:       "PIXIVFE_HOST",
		CommonName: "TCP hostname",
		// **Required**: No (ignored if PIXIVFE_UNIXSOCKET was set)
		//
		// Hostname/IP address to listen on. For example `PIXIVFE_HOST=localhost`.
	},
	{
		Name:       "PIXIVFE_PORT",
		CommonName: "TCP port",
		// **Required**: Yes (no if PIXIVFE_UNIXSOCKET was set)
		//
		// Port to listen on. For example `PIXIVFE_PORT=8745`.
	},
	{
		Name:       "PIXIVFE_UNIXSOCKET",
		CommonName: "UNIX socket path",
		// **Required**: Yes (ignored if PIXIVFE_PORT was set)
		//
		// UNIX socket to listen on. For example `PIXIVFE_UNIXSOCKET=/srv/http/pages/pixivfe`.

	},
	{
		Name:       "PIXIVFE_TOKEN",
		CommonName: "Pixiv token",
		// **Required**: Yes
		//
		// Authorization is required to fully access Pixiv's Ajax API. This variable will store your Pixiv's account cookie, which will be used by PixivFE for authorization.
		//
		// NOTE: See [How to get PIXIVFE_TOKEN](How-to-get-the-pixiv-token.md) for how to obtain your own token.

	},
	{
		Name:       "PIXIVFE_REQUESTLIMIT",
		CommonName: "limit number of request per 30 seconds",
		// **Required**: No
		//
		// Set this to a number to enable the built-in rate limiter. For example `PIXIVFE_REQUESTLIMIT=15`.
		// 
		// It might be better to enable rate limiting in the reverse proxy in front of PixivFE rather than using this.
	},
	{
		Name:       "PIXIVFE_IMAGEPROXY",
		CommonName: "image proxy server",
		Value:      BuiltinProxyUrl,
		Announce:   true,
		// **Required**: No, defaults to using the built-in proxy
		//
		// NOTE: The protocol must be included in the URL, for example `https://piximg.example.com`, where `https://` is the protocol used.
		//
		// The URL of the image proxy server. Pixiv does not allow you to fetch their images directly, requiring `Referer: https://www.pixiv.net/` to be included in the HTTP request headers. For example, trying to directly access this [image](https://i.pximg.net/img-original/img/2023/06/06/20/30/01/108783513_p0.png) returns HTTP 403 Forbidden.
		// This can be circumvented by using a reverse proxy that adds the required `Referer` HTTP request header to the HTTP request for the image. You can [host an image proxy server](Hosting-an-image-proxy-server-for-Pixiv.md), or see the [list of public image proxies](Built-in Proxy List.go). If you wish not to, or unable to get images directly from Pixiv, set this variable.
	},
	{
		Name:       "PIXIVFE_USERAGENT",
		CommonName: "user agent",
		Value:      "Mozilla/5.0 (Windows NT 10.0; rv:123.0) Gecko/20100101 Firefox/123.0",
		// **Required**: No
		//
		// The value of the `User-Agent` header, used to make requests to Pixiv's API.

	},
	{
		Name:       "PIXIVFE_ACCEPTLANGUAGE",
		CommonName: "Accept-Language header",
		Value:      "en-US,en;q=0.5",
		// **Required**: No
		//
		// The value of the `Accept-Language` header, used to make requests to Pixiv's API. You can change the response's language with this one.
	},
}

// ======================================================================
//  what lies below is irrelevant to you if you just want to use PixivFE
// ======================================================================

func CollectAllEnv() {
	for _, v := range EnvList {
		value, hasValue := os.LookupEnv(v.Name)
		if hasValue {
			v.Value = value
			v.Announce = true
		}
	}
}

func GetEnv(key string) string {
	value, _ := LookupEnv(key)
	return value
}

func LookupEnv(key string) (string, bool) {
	for _, v := range EnvList {
		if v.Name == key {
			return v.Value, v.Value != ""
		}
	}
	log.Panicf("Environment Variable Name not in `EnvironList`: %s", key)
	panic("Go's type system has no Void/noreturn type...")
}

func AnnounceAllEnv() {
	for _, v := range EnvList {
		if v.Announce {
			log.Printf("Set %s to: %s\n", v.CommonName, v.Value)
		}
	}
}
