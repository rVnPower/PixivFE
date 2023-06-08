# PixivFE

A privacy-respecting alternative front-end for Pixiv that doesn't suck

<p>
<a href="https://codeberg.org/vnpower/pixivfe">
<img alt="Get it on Codeberg" src="https://get-it-on.codeberg.org/get-it-on-blue-on-white.png" height="60">
</a>
</p>

## Features

- Lightweight - both the interface and the code
- Privacy-first - the server will do the work for you
- No bloat - we only serve HTML and CSS
- Open source - you can trust me!

## Hosting

Check out the [wiki](https://codeberg.org/VnPower/pixivfe/wiki/Home)

## Instances

| Name               | Cloudflare? | URL                      |
| ------------------ | ----------- | ------------------------ |
| exozyme (Official) | No          | https://pixivfe.exozy.me |

Hosted one yourself? Create a pull request to add it here!

## Previews

Landing page: currently borked :(

Most pages are very similar to each other. So I will only take screenshots of the artwork page and the user page, since they are the most unique.

Artwork page ([URL to this page](https://pixivfe.exozy.me/artworks/108672057)):
![Artwork](https://files.catbox.moe/c3l1vm.png)
![Artwork](https://files.catbox.moe/f9o0h0.png)
![Artwork](https://files.catbox.moe/9yc2yk.png)

User page ([URL to this page](https://pixivfe.exozy.me/users/11764388)):
![User](https://files.catbox.moe/tdc234.png)

<!-- ![Homepage](https://files.catbox.moe/053fzh.png) -->
<!-- ![Homepage](https://files.catbox.moe/m64h7s.png) -->

## To-do

- [x] Base
  - [x] Navigation bar
  - [x] Searching
  - [x] Pagination
  - [x] Configuration file
  - [x] Write a real independent API
- [ ] Index page
  - [x] Recommended artworks
  - [x] Daily rankings
  - [x] Spotlight (pixivision)
  - [ ] Discovery
    - [x] Artworks
    - [ ] Users
  - [x] Newest by all
  - [ ] Trending tags
  - [ ] Switcher (illusts/mangas)
- [ ] Single pages
  - [x] User
  - [x] Artwork
  - [ ] Spotlight
- [ ] List pages
  - [ ] Recommended artworks
  - [x] Daily rankings
  - [ ] Discovery
    - [x] Artworks
    - [ ] Users
  - [x] Newest by all
  - [x] Search results
  - [x] Switcher
- [ ] Settings
  - [ ] Login
  - [ ] Local history
  - [ ] Toggling R-18, R-18G, AI
  - [ ] Custom `pximg` proxy
- [ ] Optimization
  - [x] Split web components into smaller templates
  - [x] Clean the models + JSON
  - [x] Navigation between pages
  - [x] Lazy load images
  - [ ] Better error handling
  - [x] Fully proxy images from Pixiv

## License

[AGPL3](https://www.gnu.org/licenses/agpl-3.0.txt)
