FROM golang:1.26-alpine AS build

WORKDIR /src
COPY apps/api/go.mod ./
RUN go mod download

COPY apps/api ./
ARG APP_VERSION=dev
ARG COMMIT_SHA=unknown
RUN CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-s -w" \
  -o /out/release-sentinel-api ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /
COPY --from=build /out/release-sentinel-api /release-sentinel-api

ENV HTTP_ADDR=:8080
EXPOSE 8080

USER nonroot:nonroot
ENTRYPOINT ["/release-sentinel-api"]
