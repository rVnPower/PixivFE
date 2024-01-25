FROM docker.io/golang:1.21 as builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o pixivfe

FROM docker.io/alpine:3
COPY --from=builder /app/pixivfe /pixivfe
COPY --from=builder /app/views /views
EXPOSE 8282

ENTRYPOINT ["/pixivfe"]
