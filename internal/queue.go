package internal

import (
	"context"

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
	q.q.Push(episode, episode.GetPriority())
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

// todo: need to rename or rework that function because name is confusing
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

func worker(ctx context.Context, episode *Episode) {
	ok, err := episode.Download(ctx)

	if err != nil {
		log.WithError(err).Errorf("can't download episode %s - %s", episode.Show, episode.Title)
	}

	if ok {
		log.Infof("downloaded :: %s - %s", episode.Show, episode.Title)
	}
}
