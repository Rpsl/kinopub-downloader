##
## Build
##
FROM golang:1.18-buster AS build

WORKDIR /app

COPY . .

RUN go mod download
RUN go mod verify

RUN export CGO_ENABLED=0 && go build -ldflags "-s -w" -o /kinopub-downloader ./main.go

##
## Production
##
FROM golang:1.18-alpine

WORKDIR /app/

COPY --from=build /kinopub-downloader /app/kinopub-downloader

ENTRYPOINT ["/app/kinopub-downloader"]