package cmd

import (
	"context"

	"github.com/rpsl/kinopub-downloader/internal"
	log "github.com/sirupsen/logrus"

	"github.com/nandosousafr/podfeed"

	"github.com/rpsl/kinopub-downloader/config"
)

func Podcast(ctx context.Context) {
	cfg := ctx.Value("cfg").(*config.Config)
	queue := ctx.Value("queue").(*internal.Queue)

	for _, podcast := range cfg.Podcasts {
		pod, err := podfeed.Fetch(podcast)
		if err != nil {
			log.Error(err)
			continue
		}

		for _, ep := range pod.Items {
			episode, err := internal.NewEpisode(ep.Title, pod.Subtitle, ep.Enclosure.Url, cfg.PathForTVShows)

			if err != nil {
				log.Errorf("error processing %s - %s :: %s", pod.Subtitle, ep.Title, err)
				continue
			}

			if !episode.IsDownloaded() {
				log.Infof("marked for download :: %s - %s", pod.Subtitle, ep.Title)

				// need to move into queue implementation
				// res, err := episode.Download()
				queue.Put(episode)

				// if err != nil {
				// 	log.Error(err)
				// } else if res == true {
				// 	log.Infof("downloaded :: %s - %s", pod.Subtitle, ep.Title)
				// }
			}
		}
	}
}
