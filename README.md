# PixivFE

A privacy-respecting alternative front-end for Pixiv that doesn't suck

## What is this?

I have been researching about Pixiv's API for a long time (with a lot of problems).
One thing I noticed is that somebody actually made a decent alternative front-end for Pixiv yet
(actually, there are some, but mostly [paid + feature less](https://pixiv.moe) and/or made by Chinese and/or depends on JS).

I decided to take the lead by creating an actual front-end that is truly suckless that could access most of the features provided by Pixiv.
Intended to be written in Go with [Gin](https://gin-gonic.com).

Note that this project is still under its preparation stage.

## Previews

Homepage: currently contains recommended works, daily rankings, pixivision's articles and newest works on Pixiv. Route: `/`
![Homepage](https://files.catbox.moe/053fzh.png)
![Homepage](https://files.catbox.moe/m64h7s.png)

Artwork's page: informations about the artwork, recent artworks by the artist and related artworks. Route: `/artworks/<id>` (id for the preview: 107664519)
![Artwork](https://files.catbox.moe/x3k85p.png)
![Artwork](https://files.catbox.moe/ocy4fq.png)

## Installation

Run these commands below, then access the site on [localhost:8080](https://localhost:8080)

```
git clone https://codeberg.org/VnPower/pixivfe.git
cd pixivfe
go install main.go
go run main.go
```

## To-do

- [ ] Base
  - [ ] Navigation bar
  - [ ] Searching
  - [ ] Pagination
- [ ] Index page
  - [x] Recommended artworks
  - [x] Daily rankings
  - [x] Spotlight (pixivision)
  - [ ] Discovery
    - [ ] Artworks
    - [ ] Users
  - [x] Newest by all
  - [ ] Trending tags
  - [ ] Switcher (illusts/mangas)
- [ ] Single pages
  - [ ] User
  - [x] Artwork
  - [ ] Spotlight
- [ ] List pages
  - [ ] Recommended artworks
  - [ ] Daily rankings
  - [ ] Discovery
    - [ ] Artworks
    - [ ] Users
  - [ ] Newest by all
  - [ ] Search results
  - [ ] Switcher
- [ ] Settings
  - [ ] Login
  - [ ] Local history
  - [ ] Toggling R-18, R-18G, AI
  - [ ] Custom `pximg` proxy
- [ ] Optimization
  - [ ] Split web components into smaller templates
  - [ ] Clean the models + JSON
  - [ ] Navigation between pages

## Contributing

Every kind of contribution is appreciated! You can help the with:

- Design the front-end. If you have any interesting ideas for the front-end, you can create an issue or contact me. I desperately need you for this one ;-;
- Back-end, anything.

If you have any ideas to share, please contact me through email (vnpower@disroot.org) or [Matrix](https://matrix.to/#/@vnpower:exozy.me).
I always take ideas from everybody for this project due to it's complex nature.

## License

This project was made for learning and experimental purposes.

[~~MIT~~](https://mit-license.org/) [AGPL3](https://www.gnu.org/licenses/agpl-3.0.txt)
