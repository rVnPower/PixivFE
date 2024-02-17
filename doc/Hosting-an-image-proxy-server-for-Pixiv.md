# Hosting an i.pximg.net proxy for PixivFE

If you preferred not to use third-party image proxy server, then you could one by yourself!

To get any images from Pixiv, you just have to change the referer to Pixiv.


```
proxy_cache_path /path/to/cache levels=1:2 keys_zone=pximg:10m max_size=10g inactive=7d use_temp_path=off;

server {
    listen 443 ssl http2;

    ssl_certificate /path/to/ssl_certificate.crt;
    ssl_certificate_key /path/to/ssl_certificate.key;

    server_name pximg.example.com;
    access_log off;

    location / {
    proxy_cache pximg;
    proxy_pass https://i.pximg.net;
    proxy_cache_revalidate on;
    proxy_cache_use_stale error timeout updating http_500 http_502 http_503 http_504;
    proxy_cache_lock on;
    add_header X-Cache-Status $upstream_cache_status;
    proxy_set_header Host i.pximg.net;
    proxy_set_header Referer "https://www.pixiv.net/";

    proxy_cache_valid 200 7d;
    proxy_cache_valid 404 5m;
 }
}
```

Now, just replace `i.pximg.net` with yours, for example the image I mentioned in the environment variable page: `https://i.pximg.net/img-original/img/2023/06/06/20/30/01/108783513_p0.png` -> `https://pximg.example.com/img-original/img/2023/06/06/20/30/01/108783513_p0.png`.

You can visit this site to know more: https://pixiv.cat/reverseproxy.html. It is also an image proxy server! Try https://i.pixiv.cat/img-original/img/2023/06/06/20/30/01/108783513_p0.png.

You can also try out [this repo](https://gitler.moe/suwako/imgproxy) from TechnicalSuwako for references.