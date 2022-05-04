build-local:
	go build -o ./build/kinopub-downloader ./main.go

build-docker:
	docker build -t kinopub-downloader:local .

clean:
	rm -rf ./build

attach:
	docker run --rm -it -v $(shell pwd):/app/ --name kinopub-downloader --entrypoint 'sh' kinopub-downloader:local