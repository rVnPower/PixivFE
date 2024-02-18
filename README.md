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

| Name               | URL                          | Country | Cloudflare? | [Mozilla Observatory](https://observatory.mozilla.org/faq/) grade                                              | Uptime                                                                                                                                                                                                                                      |
| ------------------ | ---------------------------- | ------- | ----------- | -------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| exozyme (Official) | https://pixivfe.exozy.me     | US      | No          | ![Mozilla HTTP Observatory Grade](https://img.shields.io/mozilla-observatory/grade-score/pixivfe.exozy.me)     | ![Uptime Robot status](https://img.shields.io/uptimerobot/status/m796383741-c72f1ae6562dc943d032ba96) ![Uptime Robot ratio (30 days)](https://img.shields.io/uptimerobot/ratio/m796383741-c72f1ae6562dc943d032ba96?label=uptime%20%2Fmonth) |
| dragongoose        | https://pixivfe.drgns.space  | US      | No          | ![Mozilla HTTP Observatory Grade](https://img.shields.io/mozilla-observatory/grade-score/pixivfe.drgns.space)  | ![Uptime Robot status](https://img.shields.io/uptimerobot/status/m796383743-c0cf0d6b5dbb09c8dbe7dc53) ![Uptime Robot ratio (30 days)](https://img.shields.io/uptimerobot/ratio/m796383743-c0cf0d6b5dbb09c8dbe7dc53?label=uptime%20%2Fmonth) |
| ducks.party        | https://pixivfe.ducks.party  | NL      | No          | ![Mozilla HTTP Observatory Grade](https://img.shields.io/mozilla-observatory/grade-score/pixivfe.ducks.party)  | ![Uptime Robot status](https://img.shields.io/uptimerobot/status/m796383747-c92c281f520d52fe3fd894ed) ![Uptime Robot ratio (30 days)](https://img.shields.io/uptimerobot/ratio/m796383747-c92c281f520d52fe3fd894ed?label=uptime%20%2Fmonth) |
| perennialte.ch     | https://pixiv.perennialte.ch | AU      | No          | ![Mozilla HTTP Observatory Grade](https://img.shields.io/mozilla-observatory/grade-score/pixiv.perennialte.ch) | ![Uptime Robot status](https://img.shields.io/uptimerobot/status/m796383748-503799f65873a23dbc860a02) ![Uptime Robot ratio (30 days)](https://img.shields.io/uptimerobot/ratio/m796383748-503799f65873a23dbc860a02?label=uptime%20%2Fmonth) |

For information on instance uptime, see the [PixivFE instance status page](https://stats.uptimerobot.com/FbEGewWlbX).

[How to host a Pixiv image proxy](doc/Hosting-an-image-proxy-server-for-Pixiv.md)

Hosted one yourself? Create a pull request to add it here!

## License & Attributions

License: [AGPL3](https://www.gnu.org/licenses/agpl-3.0.txt)

Special thanks:

- [huggy](https://huggy.moe): author of [ugoira.com](https://ugoira.com) for the ugoira API
- [dragongoose](https://drgns.space): writing guides
- Contributors, stargazers and users like you, as well!
