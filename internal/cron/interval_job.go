package cron

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/isutare412/crawlert/internal/log"
)

type IntervalJob struct {
	name     string
	job      func()
	interval time.Duration

	wg         *sync.WaitGroup
	lifeCtx    context.Context
	lifeCancel context.CancelFunc
}

func NewIntervalJob(name string, job func(), interval time.Duration) *IntervalJob {
	lifeCtx, lifeCancel := context.WithCancel(context.Background())
	lifeCtx = log.WithValue(lifeCtx, "jobName", name)
	lifeCtx = log.WithValue(lifeCtx, "interval", interval.String())

	return &IntervalJob{
		name:       name,
		job:        job,
		interval:   interval,
		wg:         &sync.WaitGroup{},
		lifeCtx:    lifeCtx,
		lifeCancel: lifeCancel,
	}
}

func (j *IntervalJob) Run() {
	slog.InfoContext(j.lifeCtx, "run interval job")

	j.wg.Add(1)
	go func() {
		defer j.wg.Done()
		defer log.RecoverIfPanic()

		ticker := time.NewTicker(j.interval)
		defer ticker.Stop()

		for {
			j.job()

			select {
			case <-ticker.C:
			case <-j.lifeCtx.Done():
				return
			}
		}
	}()
}

func (j *IntervalJob) Shutdown() {
	j.lifeCancel()
	j.wg.Wait()

	slog.InfoContext(j.lifeCtx, "shutdwon interval job")
}
