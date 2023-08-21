# PixivFE

A privacy-respecting alternative front-end for Pixiv that doesn't suck

<p>
<a href="https://codeberg.org/vnpower/pixivfe">
<img alt="Get it on Codeberg" src="https://get-it-on.codeberg.org/get-it-on-blue-on-white.png" height="60">
</a>
</p>

Questions? Feedbacks? You can [PM me](https://matrix.to/#/@vnpower:exozy.me) on
Matrix!

You can keep track of this project's development
[here](https://codeberg.org/VnPower/pixivfe/projects/3481).

## Features

- Lightweight - both the interface and the code
- Privacy-first - the server will do the work for you
- No bloat - we only serve HTML and CSS
- Open source - you can trust me!

## Hosting

Check out [this page.](https://codeberg.org/VnPower/pixivfe/wiki/Hosting). We
currently have guides for Docker and Caddy.

Many thanks to [dragongoose](https://codeberg.org/dragongoose) for writing the
Docker guide!

## Instances

| Name               | Cloudflare? | URL                             |
| ------------------ | ----------- | ------------------------------- |
| exozyme (Official) | No          | https://pixivfe.exozy.me        |
| dragongoose        | No          | https://pixivfe.dragongoose.us  |
| WhateverItWorks    | Yes         | https://art.whateveritworks.org |

Hosted one yourself? Create a pull request to add it here!

## License

[AGPL3](https://www.gnu.org/licenses/agpl-3.0.txt)

## Note

Features like following an user, bookmarking and/or liking an artwork won't be
added anytime soon.

API routes:

- Following an user: `https://www.pixiv.net/bookmark_add.php`
- Bookmarking an artwork: `https://www.pixiv.net/ajax/illusts/bookmarks/add`

This is because for these endpoints to work, we must pass in a header called
`x-csrf-token`. The token was stored directly inside any Pixiv's pages (ex:
https://www.pixiv.net) in a variable called `pixiv.context.token`. We could
easily get this token using `curl`:

![curl https://www.pixiv.net --silent | grep pixiv.context.token](https://files.catbox.moe/pbjqtu.png)

The problem is, we cannot use this token, because when we were using `curl` to
fetch the page, we weren't authenticated. Each user has their own token,
generated every time they logout (i think). If we try to authenticate with
`PHPSESSID`, Cloudflare will stop us and we get a challenge page.

![Challenge page](https://files.catbox.moe/c1e0kp.png)

Unless we find out a way to fetch a `x-csrf-token` that works, these features
will probably never be added.

If you found out a way to fetch it, please tell me. You can test the token using
this command:
`curl -H "x-csrf-token: your_csrf_token" -X POST -d "mode=add&type=user&user_id=36055573&tag=&restrict=0&format=json" "https://www.pixiv.net/bookmark_add.php" --cookie "PHPSESSID=your_token" -A "Mozilla/5.0"`

It will return an empty JSON array if success, like this:
![Successfully followed](https://files.catbox.moe/oiwx4u.png)
