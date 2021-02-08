# Stage 1 (Build)
FROM golang:1.15-alpine3.12 AS builder

ARG VERSION
RUN apk add --update --no-cache git=2.26.2-r0 make=4.3-r0 upx=3.96-r0
WORKDIR /app/
COPY go.mod go.sum /app/
RUN go mod download
COPY . /app/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -X github.com/pterodactyl/wings/system.Version=$VERSION" \
    -v \
    -trimpath \
    -o wings \
    wings.go
RUN upx wings

# Stage 2 (Final)
FROM busybox:1.33.0
RUN echo "ID=\"busybox\"" > /etc/os-release
COPY --from=builder /app/wings /usr/bin/
CMD [ "wings", "--config", "/etc/pterodactyl/config.yml" ]
