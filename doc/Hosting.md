# Hosting PixivFE

PixivFE is a privacy-respecting alternative front-end for Pixiv that can be installed using various methods. This guide covers installation using [Docker](#docker) (recommended for production) and using a binary with a Caddy reverse proxy.

## Prerequisites

### Getting the token

PixivFE requires a Pixiv account token to access the API. Refer to the [How to get the cookie (PIXIVFE_TOKEN)](How-to-get-the-pixiv-token.md) guide for detailed instructions.

## Installation

### Docker

[Docker](https://www.docker.com/) lets you run containerized applications. Containers are loosely isolated environments that are lightweight and contain everything needed to run the application, so there's no need to rely on what's installed on the host.

There are no pre-built Docker images for PixivFE, so you'll need to build your own. PixivFE supports `amd64` and `arm64` platforms.

#### Docker Compose

Deploying PixivFE using Docker Compose requires the Compose plugin to be installed. Follow these [instructions on the Docker Docs](https://docs.docker.com/compose/install) on how to install it.

##### 1. Setting up the repository

Clone the PixivFE repository and navigate to the directory:

```bash
git clone https://codeberg.org/VnPower/PixivFE.git && cd PixivFE
```

##### 2. Set token

Copy the `PHPSESSID` cookie value into `docker/pixivfe_token.txt`. This file will be used as a [Docker secret](https://docs.docker.com/compose/use-secrets/) .

##### 3. Configure environment variables

Copy `.env.example` to `.env` and configure the variables as needed. Refer to [`Environment Variables.go`](https://codeberg.org/VnPower/PixivFE/src/branch/v2/doc/Environment%20Variables.go) for more information.

##### 4. Compose!

Run `docker compose up -d` to start PixivFE. It will be accessible at `localhost:8282`.

To view the container logs, run `docker logs -f pixivfe`.

#### Docker CLI

Deploying PixivFE using Docker CLI may be easier than Docker Compose, but requires a slightly different setup.

Furthermore, the `buildx` Docker plugin needs to be installed. Follow these [instructions on the Docker `buildx` repo](https://github.com/docker/buildx?tab=readme-ov-file#installing) on how to install it.

##### 1. Setting up the repository

Clone the PixivFE repository and navigate to the directory:

```bash
git clone https://codeberg.org/VnPower/PixivFE.git && cd PixivFE
```

##### 2. Building the image

For `amd64` platforms:

```bash
docker buildx build --platform linux/amd64 -t vnpower/pixivfe:latest --load .
```

For `arm64` platforms:

```bash
docker buildx build --platform linux/arm64 -t vnpower/pixivfe:latest-arm64 --load .
```

##### 3. Configure environment variables

Copy `.env.example` to `.env` and configure the variables as needed. Refer to [`Environment Variables.go`](https://codeberg.org/VnPower/PixivFE/src/branch/v2/doc/Environment%20Variables.go) for more information.

##### 4. Deploying PixivFE


Run the following command to deploy PixivFE:

> **Note**:
>
> If using an `arm64` platform, use the `vnpower/pixivfe:latest-arm64` image.

```bash
docker run -d --name pixivfe -p 8282:8282 --env-file .env vnpower/pixivfe:latest
```

To use a different host port (e.g., 8080):

```bash 
docker run -d --name pixivfe -p 8080:8282 --env-file .env vnpower/pixivfe:latest
```

When using a reverse proxy, prefix the host port with `127.0.0.1` to make PixivFE listen only on the local host port (e.g., `127.0.0.1:8080:8282`).

### Binary with Caddy reverse proxy

##### 1. Setting up the repository

Clone the PixivFE repository, navigate to the directory, and install the dependencies:

```bash
git clone https://codeberg.org/VnPower/PixivFE.git && cd PixivFE
go install
```

##### 2. Configure environment variables

Copy `.env.example` to `.env` and configure the variables as needed. Refer to [`Environment Variables.go`](https://codeberg.org/VnPower/PixivFE/src/branch/v2/doc/Environment%20Variables.go) for more information.

##### 3. Deploying PixivFE

Run `env $(cat .env | xargs) go run main.go` to start PixivFE. It will be accessible at `localhost:8282`.

##### 4. Deploying Caddy

[Caddy](https://caddyserver.com/) is a great alternative to NGINX, because it is written in Go but also easy to config.

Install Caddy using your package manager.

In the PixivFE directory, create a `Caddyfile` with the following content:

```caddy
example.com {
  reverse_proxy localhost:8282
}
```

Replace `example.com` with your domain and `8282` with the PixivFE port if you changed it.

Run `caddy run` to start Caddy.

## Acknowledgement

- [Keep Caddy Running](https://caddyserver.com/docs/running#keep-caddy-running)
