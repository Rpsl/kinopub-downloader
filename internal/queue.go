package internal

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/oleiade/lane"
	"github.com/remeh/sizedwaitgroup"

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
	var swg = sizedwaitgroup.New(2)

	for {
		select {
		case <-ctx.Done():
			log.Infoln("all workers finished")
			return
		default:
			log.Debugf("waiting for queueu")

			for q.q.Size() > 0 {
				log.Debugf("queue have %d items", q.q.Size())
				item, _ := q.q.Pop()
				ep := item.(*Episode)

				swg.Add()

				go worker(ep, &swg)
			}
			swg.Wait()

			time.Sleep(time.Second)
		}
	}
}

func parsePriority(episode *Episode) int64 {
	priority, err := strconv.ParseInt(fmt.Sprintf("%d0%d", episode.SeasonNumber, episode.EpisodeNumber), 10, 64)

	log.Debugf("%s - priority is: %d", episode.Title, priority)

	if err != nil {
		log.Warnf("can't parse priority for queue from %s - %s", episode.TVShow, episode.Title)

		priority = 999
	}

	return priority
}

func worker(episode *Episode, swg *sizedwaitgroup.SizedWaitGroup) {
	ok, err := episode.Download()

	if err != nil {
		log.WithError(err).Errorf("can't download episode %s - %s", episode.TVShow, episode.Title)
	}

	if ok {
		log.Infof("downloaded :: %s - %s", episode.TVShow, episode.Title)
	}

	swg.Done()
}
