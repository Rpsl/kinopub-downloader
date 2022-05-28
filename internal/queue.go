package internal

import (
	"context"
	"fmt"
	"strconv"

	"github.com/oleiade/lane"
	"github.com/remeh/sizedwaitgroup"
	"github.com/rpsl/kinopub-downloader/config"

	log "github.com/sirupsen/logrus"
)

type Queue struct {
	q       *lane.PQueue
	channel chan Episode
}

func NewQueue() *Queue {
	var priorityQueue = lane.NewPQueue(lane.MINPQ)

	return &Queue{
		q:       priorityQueue,
		channel: make(chan Episode),
	}
}

func (q *Queue) Put(episode *Episode) {
	priority := parsePriority(episode)

	q.q.Push(episode, int(priority))
}

func (q *Queue) Work(ctx context.Context) {
	cfg := ctx.Value(CtxCfgKey).(*config.Config)
	swg := sizedwaitgroup.New(cfg.ConcurrencyDownloads)

	qChannel := make(chan interface{}, cfg.ConcurrencyDownloads)

	go q.fillQueue(ctx, qChannel)

	for {
		select {
		case <-ctx.Done():
			// need to wait workers for graceful shutdown and deleting  temporary files
			swg.Wait()
			log.Infoln("all workers finished")

			return
		case item := <-qChannel:
			ep := item.(*Episode)

			swg.Add()

			go func() {
				log.Debugln("Worker is started")

				defer swg.Done()

				worker(ctx, ep)
				log.Debugln("Worker finished")
			}()
		}
	}
}

func (q *Queue) fillQueue(ctx context.Context, qChannel chan interface{}) {
	cfg := ctx.Value(CtxCfgKey).(*config.Config)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if q.q.Size() > 0 && len(qChannel) < cfg.ConcurrencyDownloads {
				log.Debugf("queue have %d items", q.q.Size())

				item, _ := q.q.Pop()
				qChannel <- item
			}
		}
	}
}

func parsePriority(episode *Episode) int64 {
	priority, err := strconv.ParseInt(fmt.Sprintf("%d0%d", episode.SeasonNumber, episode.EpisodeNumber), 10, 64)

	if err != nil {
		log.Warnf("can't parse priority for queue from %s - %s", episode.Show, episode.Title)

		priority = 999
	}

	return priority
}

func worker(ctx context.Context, episode *Episode) {
	ok, err := episode.Download(ctx)

	if err != nil {
		log.WithError(err).Errorf("can't download episode %s - %s", episode.Show, episode.Title)
	}

	if ok {
		log.Infof("downloaded :: %s - %s", episode.Show, episode.Title)
	}
}
