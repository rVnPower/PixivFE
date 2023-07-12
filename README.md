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

With [Docker](https://codeberg.org/VnPower/pixivfe/wiki/Hosting-with-Docker) or with [Caddy](https://codeberg.org/VnPower/pixivfe/wiki/Hosting-with-Caddy).

Many thanks to [dragongoose](https://codeberg.org/dragongoose) for writing the Docker page!

## Instances

| Name               | Cloudflare? | URL                            |
|--------------------|-------------|--------------------------------|
| exozyme (Official) | No          | https://pixivfe.exozy.me       |
| dragongoose        | No          | https://pixivfe.dragongoose.us |

Hosted one yourself? Create a pull request to add it here!

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
  - [x] Login
  - [ ] Local history
  - [ ] Toggling R-18, R-18G, AI (?)
  - [x] Custom `pximg` proxy
- [ ] Optimization
  - [x] Split web components into smaller templates
  - [x] Clean the models + JSON
  - [x] Navigation between pages
  - [x] Lazy load images
  - [x] Better error handling
  - [x] Fully proxy images from Pixiv
  - [ ] Optimize pagination code

## License

[AGPL3](https://www.gnu.org/licenses/agpl-3.0.txt)
