# PixivFE

A privacy-respecting alternative front-end for Pixiv that doesn't suck.

<p>
<a href="https://codeberg.org/vnpower/pixivfe">
<img alt="Get it on Codeberg" src="https://get-it-on.codeberg.org/get-it-on-blue-on-white.png" height="60">
</a>
</p>

![CI badge](https://ci.codeberg.org/api/badges/12556/status.svg)
[![Go Report Card](https://goreportcard.com/badge/codeberg.org/vnpower/pixivfe/v2)](https://goreportcard.com/report/codeberg.org/vnpower/pixivfe)

Questions? Feedback? You can [PM me](https://matrix.to/#/@vnpower:eientei.org) on Matrix! You can also see the [Known quirks](doc/Quirks.md) page to check if your issue has a known solution.

You can keep track of this project's development using the [roadmap](doc/dev/general.md).

## Features

- Lightweight - both the interface and the code
- Privacy-first - the server will do the work for you
- No bloat - we only serve HTML, CSS and minimal JS code
- Open source - you can trust me!

## Hosting

You can use PixivFE for personal use! Assuming that you use an operating system that can run POSIX shell scripts, install `go`, clone this repository, modify the `run.sh` file, and profit!
I recommend self-hosting your own instance for personal use, instead of relying entirely on official instances.

To deploy PixivFE using Docker or the compiled binary, see the [Hosting PixivFE](doc/Hosting.md) wiki page.

PixivFE can work with or without an external image proxy server. Here is [the built-in proxy list](doc/Built-in%20Proxy%20List.go).
See [hosting a Pixiv image proxy](doc/Hosting-an-image-proxy-server-for-Pixiv.md) if you want to host one yourself.


## Development

**Requirements:**

- [Go](https://go.dev/doc/install) (to build PixivFE from source)
- [Sass](https://github.com/sass/dart-sass/) (will be run by PixivFE in development mode)

To install Dart Sass, you can choose any of the following methods.

- use system package manager (usually called `dart-sass`)
- download executable from [the official release page](https://github.com/sass/dart-sass/releases)
- `pnpm i -g sass`

```bash
# Clone the PixivFE repository
git clone https://codeberg.org/VnPower/PixivFE.git && cd PixivFE

# Run in PixivFE in development mode (styles and templates reload automatically)
PIXIVFE_DEV=1 <other_environment_variables> go run .
```

## Instances

<!-- The current instance table is really wide; maybe there's a better way of formatting it without losing information?
The badges are also difficult to read on a small screen due to Codeberg shrinking the width of the columns -->

| Name               | URL                          | Country | Cloudflare? | [Observatory](https://observatory.mozilla.org/faq/) grade                                                                                                              | Uptime                                                                                                                                                                                                                                      |
| ------------------ | ---------------------------- | ------- | ----------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| exozyme (Official) | https://pixivfe.exozy.me     | US      | No          | [![Mozilla HTTP Observatory Grade](https://img.shields.io/mozilla-observatory/grade-score/pixivfe.exozy.me?label=)](https://observatory.mozilla.org/analyze/pixivfe.exozy.me)         | ![Uptime Robot status](https://img.shields.io/uptimerobot/status/m796383741-c72f1ae6562dc943d032ba96?&cacheSeconds=3600) ![Uptime Robot ratio (30 days)](https://img.shields.io/uptimerobot/ratio/m796383741-c72f1ae6562dc943d032ba96?label=uptime%20%2Fmonth&cacheSeconds=3600) |
| dragongoose        | https://pixivfe.drgns.space  | US      | No          | [![Mozilla HTTP Observatory Grade](https://img.shields.io/mozilla-observatory/grade-score/pixivfe.drgns.space?label=)](https://observatory.mozilla.org/analyze/pixivfe.drgns.space)   | ![Uptime Robot status](https://img.shields.io/uptimerobot/status/m796383743-c0cf0d6b5dbb09c8dbe7dc53?&cacheSeconds=3600) ![Uptime Robot ratio (30 days)](https://img.shields.io/uptimerobot/ratio/m796383743-c0cf0d6b5dbb09c8dbe7dc53?label=uptime%20%2Fmonth&cacheSeconds=3600) |
| ducks.party        | https://pixivfe.ducks.party  | NL      | No          | [![Mozilla HTTP Observatory Grade](https://img.shields.io/mozilla-observatory/grade-score/pixivfe.ducks.party?label=)](https://observatory.mozilla.org/analyze/pixivfe.ducks.party)   | ![Uptime Robot status](https://img.shields.io/uptimerobot/status/m796383747-c92c281f520d52fe3fd894ed?&cacheSeconds=3600) ![Uptime Robot ratio (30 days)](https://img.shields.io/uptimerobot/ratio/m796383747-c92c281f520d52fe3fd894ed?label=uptime%20%2Fmonth&cacheSeconds=3600) |
| perennialte.ch     | https://pixiv.perennialte.ch | AU      | No          | [![Mozilla HTTP Observatory Grade](https://img.shields.io/mozilla-observatory/grade-score/pixiv.perennialte.ch?label=)](https://observatory.mozilla.org/analyze/pixiv.perennialte.ch) | ![Uptime Robot status](https://img.shields.io/uptimerobot/status/m796383748-503799f65873a23dbc860a02?&cacheSeconds=3600) ![Uptime Robot ratio (30 days)](https://img.shields.io/uptimerobot/ratio/m796383748-503799f65873a23dbc860a02?label=uptime%20%2Fmonth&cacheSeconds=3600) |

If you are hosting your own instance, you can create a pull request to add it here!

For more information on instance uptime, see the [PixivFE instance status page](https://stats.uptimerobot.com/FbEGewWlbX).

## License

License: [AGPL3](https://www.gnu.org/licenses/agpl-3.0.txt)
