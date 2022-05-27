package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	cli "github.com/jawher/mow.cli"
	"github.com/rpsl/kinopub-downloader/cmd"
	"github.com/rpsl/kinopub-downloader/config"
	"github.com/rpsl/kinopub-downloader/internal"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: time.RFC3339,
		FullTimestamp:   true,
	})
	// log.SetLevel(log.DebugLevel)

	cfg, err := config.LoadConfig()

	if err != nil {
		log.WithError(err).Fatal("failed to load configuration file")
	}

	var wg sync.WaitGroup
	wg.Add(2)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "cfg", cfg)
	ctx = context.WithValue(ctx, "wg", &wg)
	ctx = context.WithValue(ctx, "queue", internal.NewQueue())
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGTERM)
	defer stop()

	go runCli(ctx)
	go runQueue(ctx)

	wg.Wait()
	fmt.Println("Done!")
}

func runCli(ctx context.Context) {
	wg := ctx.Value("wg").(*sync.WaitGroup)
	defer wg.Done()

	app := cli.App("Kinopub Downloader", "Utility for downloading Movies & TV Shows from kino.pub (work only with pro account)")

	app.Command("podcast", "download by podcasts feed", func(c *cli.Cmd) {
		c.Action = func() {
			cmd.Podcast(ctx)
		}
	})

	app.Run(os.Args)
}

func runQueue(ctx context.Context) {
	wg := ctx.Value("wg").(*sync.WaitGroup)
	defer wg.Done()

	log.Infoln("starting queue")
	queue := ctx.Value("queue").(*internal.Queue)

	// todo will run endless because runCLI run only once

	select {
	case <-ctx.Done():
		log.Infoln("all workers finished")
	default:
		queue.Work()
	}

}
