FROM golang:1.23-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/daemon ./cmd/daemon
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/migrate ./cmd/migrate

FROM alpine:3.21

RUN apk add --no-cache ca-certificates \
    && addgroup -S kadeem \
    && adduser -S -D -H -h /app -s /sbin/nologin -G kadeem kadeem

WORKDIR /app

COPY --from=build /out/daemon /usr/local/bin/daemon
COPY --from=build /out/migrate /usr/local/bin/migrate
COPY --chown=kadeem:kadeem migrations ./migrations

USER kadeem

EXPOSE 8080

CMD ["daemon"]
