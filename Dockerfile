# syntax=docker/dockerfile:1.7

ARG GO_VERSION=1.26.3

FROM golang:${GO_VERSION}-alpine AS base
WORKDIR /src
RUN apk add --no-cache ca-certificates git

FROM base AS deps
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

FROM deps AS test
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go test ./...

FROM deps AS build
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/kappelas-go-bot .

FROM alpine:3.22 AS runtime
RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S -G app -h /home/app app \
    && mkdir -p /home/app/.config/kappelas-go-bot/data \
    && chown -R app:app /home/app

USER app
WORKDIR /app
ENV GIN_MODE=release \
    SERVER_PORT=8080 \
    HOME=/home/app

COPY --from=build /out/kappelas-go-bot /app/kappelas-go-bot

EXPOSE 8080
ENTRYPOINT ["/app/kappelas-go-bot"]
