# PixivFE

A privacy-respecting alternative front-end for Pixiv that doesn't suck.

<p>
<a href="https://codeberg.org/vnpower/pixivfe">
<img alt="Get it on Codeberg" src="https://get-it-on.codeberg.org/get-it-on-blue-on-white.png" height="60">
</a>
</p>

![CI badge](https://ci.codeberg.org/api/badges/12556/status.svg)
[![Go Report Card](https://goreportcard.com/badge/codeberg.org/vnpower/pixivfe)](https://goreportcard.com/report/codeberg.org/vnpower/pixivfe)

Questions? Feedback? You can [PM me](https://matrix.to/#/@vnpower:eientei.org) on Matrix! You can also look in [Known Quirks Of PixivFE](doc/Quirks.md) to see if your issue already has a known solution.

You can keep track of this project's development [here](doc/dev/Things-to-do.md).

## Features

- Lightweight - both the interface and the code
- Privacy-first - the server will do the work for you
- No bloat - we only serve HTML, CSS and minimal JS code
- Open source - you can trust me!

## Hosting

You can use PixivFE for personal use! Assuming that you use an operating system that can run POSIX shell scripts, install `go`, clone this repository, modify the `run.sh` file, and profit!
I recommend self-hosting your own instance for personal use, instead of relying entirely on official instances.

[How to host PixivFE using Docker, or Caddy](doc/Hosting.md)

## Development

**Requirements:**

- [pnpm](https://pnpm.io/installation) (to install Sass)
- [go](https://go.dev/doc/install) (to build PixivFE from source)

```bash
# Clone the PixivFE repository
git clone https://codeberg.org/VnPower/PixivFE.git && cd PixivFE

# Install Sass globally using pnpm
pnpm i -g sass

# Compile styles using Sass and watch for changes
sass --watch ./views/css

# Run in PixivFE in development mode (templates reload automatically)
PIXIVFE_DEV=1 <other_environment_variables> go run .
```

## Instances

| Name               | Cloudflare? | URL                             |
| ------------------ | ----------- | ------------------------------- |
| exozyme (Official) | No          | https://pixivfe.exozy.me        |
| dragongoose        | No          | https://pixivfe.drgns.space     |
| chaotic.ninja      | No          | https://pix.chaotic.ninja       |
| WhateverItWorks    | Yes         | https://art.whateveritworks.org |
| ducks.party        | No          | https://pixivfe.ducks.party     |
| perennialte.ch     | No          | https://pixiv.perennialte.ch    |

[How to host a Pixiv image proxy](doc/Hosting-an-image-proxy-server-for-Pixiv.md)

Hosted one yourself? Create a pull request to add it here!

## License & Attributions

License: [AGPL3](https://www.gnu.org/licenses/agpl-3.0.txt)

Special thanks:

- [huggy](https://huggy.moe): author of [ugoira.com](https://ugoira.com) for the ugoira API
- [dragongoose](https://drgns.space): writing guides
- Contributors, stargazers and users like you, as well!
