# ------ Builder stage ------
FROM docker.io/golang:1.22 as builder
WORKDIR /app

COPY go.* ./
RUN go mod download
COPY . ./

# Build the application binary with optimisations for a smaller, static binary
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -ldflags="-s -w" -o pixivfe

# ------ Final image ------
FROM docker.io/alpine:3.19
WORKDIR /app

# Create a non-root user `pixivfe` for security purposes and set ownership
RUN addgroup -g 1000 -S pixivfe && \
    adduser -u 1000 -S pixivfe -G pixivfe && \
    chown -R pixivfe:pixivfe /app

# Copy the compiled application and other necessary files from the builder stage
COPY --from=builder /app/pixivfe /app/pixivfe
COPY --from=builder /app/views /app/views
COPY ./docker/entrypoint.sh /entrypoint.sh
# Include entrypoint script and ensure it's executable
RUN chmod +x /entrypoint.sh && \
    chown pixivfe:pixivfe /entrypoint.sh

# Use the non-root user to run the application
USER pixivfe

EXPOSE 8282

ENTRYPOINT ["/entrypoint.sh"]

HEALTHCHECK --interval=30s --timeout=3s --start-period=15s --start-interval=5s --retries=3 \
 CMD wget --spider -q --tries=1 http://127.0.0.1:8282/about || exit 1
