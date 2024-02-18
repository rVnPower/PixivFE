#!/bin/sh

# Update the program every time you run?
# git pull

# Visit ./doc/Environment\ Variables.go for more details
export PIXIVFE_TOKEN=token_123456
export PIXIVFE_IMAGEPROXY=pximg.cocomi.cf
# export PIXIVFE_UNIXSOCKET=/srv/http/pages/pixivfe
export PIXIVFE_PORT=8282

go mod download
go get codeberg.org/vnpower/pixivfe/v2/...
CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o pixivfe

./pixivfe
