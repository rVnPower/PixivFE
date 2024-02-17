# Environment Variables

PixivFE's behavior is governed by those Environment Variables.

### PIXIVFE_TOKEN
**Required**: Yes

Authorization is required to fully access Pixiv's Ajax API. This variable will store your Pixiv's account cookie, which will be used by PixivFE for authorization. 

**Notice:** Please read [How to get PIXIVFE_TOKEN](How-to-get-the-cookie-(PIXIVFE_TOKEN).md) to see how can you get your own token and more.

### PIXIVFE_PORT
**Required**: Yes (no if PIXIVFE_UNIXSOCKET was set)

Port to run on. For example `PIXIVFE_PORT=8745`. 

### PIXIVFE_UNIXSOCKET
**Required**: Yes (ignored if PIXIVFE_PORT was set)

UNIX socket to run on. For example `PIXIVFE_UNIXSOCKET=/srv/http/pages/pixivfe`. 

### PIXIVFE_IMAGEPROXY
**Required**: Yes

See the current [list of image proxies](https://pixivfe.exozy.me/settings).

The address to proxy images. Pixiv does not allow you to get their images normally. For example, this [image](https://i.pximg.net/img-original/img/2023/06/06/20/30/01/108783513_p0.png). We could bypass this anyway by using NGINX and reverse proxy. [You can host an image proxy server if you want](./Hosting-an-image-proxy-server-for-Pixiv.md). If you wish not to, or unable to get images directly from Pixiv, set this variable. 

### PIXIVFE_USERAGENT
**Required**: No

Default: Mozilla/5.0

The value of the `User-Agent` header, used to make requests to Pixiv's API.

### PIXIVFE_BASEURL
**Required**: No

Used to generate meta tags.

### PIXIVFE_ACCEPTLANGUAGE
**Required**: No
Default: en-US,en;q=0.5

The value of the `Accept-Language` header, used to make requests to Pixiv's API. You can change the response's language with this one.