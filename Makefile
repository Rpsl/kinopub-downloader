build-bin:
	go build -ldflags "-s -w" -o ./build/kinopub-downloader ./main.go

build-docker:
	docker build -t kinopub-downloader:local .

build-clean:
	rm -rf ./build

attach:
	docker run --rm -it -v $(shell pwd):/app/ --name kinopub-downloader --entrypoint 'sh' kinopub-downloader:local

test:
	go test -v ./...

lint:
	golangci-lint run --out-format=github-actions

mod-update:
	go get -u && go mod tidy
