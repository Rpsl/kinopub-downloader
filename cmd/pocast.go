package cmd

import (
	"context"
	"time"

	"github.com/rpsl/kinopub-downloader/internal"
	log "github.com/sirupsen/logrus"

	"github.com/nandosousafr/podfeed"
	"github.com/rpsl/kinopub-downloader/config"
)

func Podcast(ctx context.Context) {
	cfg := ctx.Value(internal.CtxCfgKey).(*config.Config)
	queue := ctx.Value(internal.CtxQueueKey).(*internal.Queue)

	// for first update on starting program
	updateFeeds(cfg, queue)

	timer := time.NewTimer(time.Hour * time.Duration(cfg.HoursToRefresh))
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			if !timer.Stop() {
				updateFeeds(cfg, queue)
				timer.Reset(time.Hour * time.Duration(cfg.HoursToRefresh))
			}
		}
	}
}

func updateFeeds(cfg *config.Config, queue *internal.Queue) {
	log.Debugf("star updating all feeds")

	for _, podcast := range cfg.Podcasts {
		pod, err := podfeed.Fetch(podcast)
		if err != nil {
			log.Error(err)
			continue
		}

		for _, ep := range pod.Items {

			// todo extract title's routines in separate functions
			show := ""

			switch {
			case pod.Subtitle != "":
				show = pod.Subtitle
			case pod.Title != "":
				show = pod.Title
			default:
				log.Errorf("can't detech show name for %s", podcast)
				continue
			}

			episode, err := internal.NewEpisode(ep.Title, show, ep.Enclosure.Url, config.PathForTVShows)

			if err != nil {
				log.Errorf("error processing %s - %s :: %s", pod.Subtitle, ep.Title, err)
				continue
			}

			if !episode.IsDownloaded() {
				log.Infof("marked for download :: %s - %s", pod.Subtitle, ep.Title)

				queue.Put(episode)
			}
		}
	}
}
