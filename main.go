package main

import (
	"os"
	"time"

	"github.com/rpsl/kinopub-downloader/cmd"
	log "github.com/sirupsen/logrus"

	cli "github.com/jawher/mow.cli"

	"github.com/rpsl/kinopub-downloader/config"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: time.RFC3339,
		FullTimestamp:   true,
	})

	cfg, err := config.LoadConfig()

	if err != nil {
		log.WithError(err).Fatal("failed to load configuration file")
	}

	app := cli.App("Kinopub Downloader", "Utility for downloading Movies & TV Shows from kino.pub (work only with pro account)")

	app.Command("podcast", "download by podcasts feed", func(c *cli.Cmd) {
		c.Action = func() {
			cmd.Podcast(cfg)
		}
	})

	app.Run(os.Args)
}
