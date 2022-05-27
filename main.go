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
	log.SetLevel(log.DebugLevel)

	cfg, err := config.LoadConfig()

	if err != nil {
		log.WithError(err).Fatal("failed to load configuration file")
	}

	var wg sync.WaitGroup

	wg.Add(2)

	ctx := context.Background()
	ctx = context.WithValue(ctx, internal.CtxCfgKey, cfg)
	ctx = context.WithValue(ctx, internal.CtxWgKey, &wg)
	ctx = context.WithValue(ctx, internal.CtxQueueKey, internal.NewQueue())
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)

	defer stop()

	go runCli(ctx)
	go runQueue(ctx)

	wg.Wait()

	<-ctx.Done()
	fmt.Println("shutting down gracefully, press Ctrl+C again to force")
}

func runCli(ctx context.Context) {
	wg := ctx.Value(internal.CtxWgKey).(*sync.WaitGroup)
	defer wg.Done()

	app := cli.App("Kinopub Downloader", "Utility for downloading Movies & TV Shows from kino.pub (work only with pro account)")

	app.Command("podcast", "download by podcasts feed", func(c *cli.Cmd) {
		c.Action = func() {
			cmd.Podcast(ctx)
		}
	})

	app.Run(os.Args)
	log.Debugf("runCLI finished")
}

func runQueue(ctx context.Context) {
	wg := ctx.Value(internal.CtxWgKey).(*sync.WaitGroup)
	defer wg.Done()

	log.Infoln("starting queue")

	queue := ctx.Value(internal.CtxQueueKey).(*internal.Queue)
	queue.Work(ctx)
	log.Debugf("runQueue finished")
}
