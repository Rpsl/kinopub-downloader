# Kino.pub downloader

It's utility for downloading media content from [Kino.Pub](https://kino.pub). You need PRO account for using. At now
it's support only downloading by podcasts links

## How to use

### Build docker image

```shell
$ make build-docker
``` 

### Create config file

```shell
$ cp config.example.toml config.toml
```

### Run docker image

```shell
docker run --rm -v /path/to/config.toml:/app/config.toml -v /path/to/download/folder:/app/data --name kinopub-downloader kinopub-downloader:local 
```

### Proxy

Because [Kino.Pub](https://kino.pub) was ban in some countries you need use vpn or proxy if you have problems with
download. You can set up proxy server as env variable

```shell
$ export HTTPS_PROXY="(http|socks5)://proxyIp:proxyPort"
# also it work with docker image
$ docker run .... -e HTTPS_PROXY="(http|socks5)://proxyIp:proxyPort" ....
```
