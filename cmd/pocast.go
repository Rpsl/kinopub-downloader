package cmd

import (
	"github.com/rpsl/kinopub-downloader/internal"
	log "github.com/sirupsen/logrus"

	"github.com/nandosousafr/podfeed"

	"github.com/rpsl/kinopub-downloader/config"
)

func Podcast(config *config.Config) {
	for _, podcast := range config.Podcasts {
		pod, err := podfeed.Fetch(podcast)
		if err != nil {
			log.Error(err)
			continue
		}

		for _, ep := range pod.Items {
			episode, err := internal.NewEpisode(ep.Title, pod.Subtitle, ep.Enclosure.Url, config.PathForTVShows)

			if err != nil {
				log.Errorf("error processing %s - %s :: %s", pod.Subtitle, ep.Title, err)
				continue
			}

			if episode.IsDownloaded() == false {
				log.Infof("marked for download :: %s - %s", pod.Subtitle, ep.Title)

				// need to move into queue implementation
				res, err := episode.Download()

				if err != nil {
					log.Error(err)
				} else if res == true {
					log.Infof("downloaded :: %s - %s", pod.Subtitle, ep.Title)
				}
			}
		}
	}
}