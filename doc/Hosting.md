# Hosting PixivFE

This page covers multiple methods to install PixivFE. Using [Docker](#docker) is recommended for production use.

## Prerequisites

### Getting the token

PixivFE needs an account token to reach the API.

You can check out [this page](How-to-get-the-pixiv-token.md) for detailed information about how to get the token.

## Installation

### Docker

Docker images for PixivFE can be built with support for `amd64` and `arm64` platforms.

However, there is no Docker image for PixivFE, so you will have to build your own.

#### Docker Compose

Deploying PixivFE using Docker Compose requires the Compose plugin to be installed. Follow these [instructions on the Docker Docs](https://docs.docker.com/compose/install) on how to install it.

##### 1. Setting up the repository

Clone the repo and `cd` into the directory:

```bash
git clone https://codeberg.org/VnPower/PixivFE.git && cd PixivFE
```

##### 2. Set token

A [secret](https://docs.docker.com/compose/use-secrets/) is used to provide the token used by PixivFE to fetch content.

Copy the contents of the `PHPSESSID` cookie into `docker/pixivfe_token.txt`.

##### 3. Compose!

```bash
docker compose up -d
```

Your PixivFE instance is now up at `localhost:8282`!

To follow container logs:

```bash
docker logs -f pixivfe
```

#### Docker CLI

Deploying PixivFE using Docker CLI may be easier than Docker Compose, but requires a slightly different setup.

Furthermore, the `buildx` Docker plugin needs to be installed. Follow these [instructions on the Docker `buildx` repo](https://github.com/docker/buildx?tab=readme-ov-file#installing) on how to install it.

##### 1. Setting up the repository

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

##### 3. Deploying PixivFE

Deploy PixivFE:

```
docker run -d --name pixivfe -p 8282:8282 vnpower/pixivfe:latest
```

Deploy using a different port on the host (in this case, port 8080):

```
docker run -d --name pixivfe -p 8080:8282 vnpower/pixivfe:latest
```

> **Note**:
>
> If deploying on an `arm64` platform, use the `vnpower/pixivfe:latest-arm64` image instead.

If you're using a reverse proxy in front of PixivFE, prefix the port numbers with `127.0.0.1` so that PixivFE only listens on the host port **locally**. For example, if the host port for PixivFE is `8080`, specify `127.0.0.1:8080:8282`. 

### Binary with Caddy reverse proxy

Clone the repository and install the dependencies.

```bash
git clone https://codeberg.org/VnPower/PixivFE.git && cd PixivFE
go install
```

You may wanted to check out some of the environment variables used by PixivFE before continuing.

After that, run `go run main.go`. And PixivFE should be running now!

[Caddy](https://caddyserver.com/) is a great alternative to NGINX, because it is written in Go but also easy to config.

Install Caddy using your package manager.

After installing Caddy, make sure that you are inside PixivFE's directory. Then, create a file named `Caddyfile`. You should see something like this:

```
example.com {
  reverse_proxy localhost:8282
}
```

Change `example.com` to your domain, also change `8282` if you set the PixivFE's port to something else.

Finally, run `caddy run`.

## Acknowledgement

- [Keep Caddy Running](https://caddyserver.com/docs/running#keep-caddy-running)
