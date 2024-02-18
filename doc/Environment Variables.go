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
	Value      string // available at run-time
	CommonName string
}

// All environment variables used by PixivFE
var EnvironList []EnvVar = []EnvVar{
	{
		Name:       "PIXIVFE_DEV",
		CommonName: "development mode",
		// **Required**: No
		//
		// Set this to anything to enable development mode, in which the server will live-reload HTML templates and disable caching.
	},

	{
		Name:       "PIXIVFE_HOST",
		CommonName: "HOST",
		// **Required**: Yes (no if PIXIVFE_UNIXSOCKET was set)
		//
		// Hostname/IP address to listen on. For example `PIXIVFE_HOST=localhost`.
	},
	{
		Name:       "PIXIVFE_PORT",
		CommonName: "PORT",
		// **Required**: Yes (no if PIXIVFE_UNIXSOCKET was set)
		//
		// Port to listen on. For example `PIXIVFE_PORT=8745`.
	},
	{
		Name:       "PIXIVFE_UNIXSOCKET",
		CommonName: "UNIXSOCKET",
		// **Required**: Yes (ignored if PIXIVFE_PORT was set)
		//
		// UNIX socket to listen on. For example `PIXIVFE_UNIXSOCKET=/srv/http/pages/pixivfe`.

	},
	{
		Name:       "PIXIVFE_TOKEN",
		CommonName: "TOKEN",
		// **Required**: Yes
		//
		// Authorization is required to fully access Pixiv's Ajax API. This variable will store your Pixiv's account cookie, which will be used by PixivFE for authorization.
		//
		// **Notice:** Please read [How to get PIXIVFE_TOKEN](How-to-get-the-pixiv-token.md) to see how can you get your own token and more.

	},
	{
		Name:       "PIXIVFE_IMAGEPROXY",
		CommonName: "IMAGEPROXY",
		// **Required**: Yes
		//
		// See the current [list of image proxies](Built-in Proxy List.go).
		//
		// The address to proxy images. Pixiv does not allow you to get their images normally. For example, this [image](https://i.pximg.net/img-original/img/2023/06/06/20/30/01/108783513_p0.png). We could bypass this anyway by using NGINX and reverse proxy. [You can host an image proxy server if you want](./Hosting-an-image-proxy-server-for-Pixiv.md). If you wish not to, or unable to get images directly from Pixiv, set this variable.
	},
	{
		Name:       "PIXIVFE_REQUESTLIMIT",
		CommonName: "REQUESTLIMIT",
		Value:      "15",
		// **Required**: No
	},
	{
		Name:       "PIXIVFE_USERAGENT",
		CommonName: "USERAGENT",
		Value:      "Mozilla/5.0",
		// **Required**: No
		//
		// The value of the `User-Agent` header, used to make requests to Pixiv's API.

	},
	{
		Name:       "PIXIVFE_ACCEPTLANGUAGE",
		CommonName: "ACCEPTLANGUAGE",
		Value:      "en-US,en;q=0.5",
		// **Required**: No
		//
		// The value of the `Accept-Language` header, used to make requests to Pixiv's API. You can change the response's language with this one.
	},
}

func CollectEnv() {
	for _, v := range EnvironList {
		value, hasValue := os.LookupEnv(v.Name)
		if hasValue {
			v.Value = value
		}
	}
}

func LookupEnv(key string) (string, bool) {
	for _, v := range EnvironList {
		if v.Name == key {
			return v.Value, v.Value == ""
		}
	}
	log.Panicf("Environment Variable Name not in `EnvironList`: %s", key)
	panic("Go's type system has no Void/noreturn type...")
}
