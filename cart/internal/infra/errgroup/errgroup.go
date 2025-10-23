package errgroup

import (
	"sync"
	"time"

	"golang.org/x/net/context"
)

type Group struct {
	ctx      context.Context
	wg       *sync.WaitGroup
	limit    int
	duration time.Duration
	chIn     chan func() error
	chOut    chan error
}

func WithContext(ctx context.Context, limit, bufferSize int, duration time.Duration) (*Group, context.Context) {
	chIn := make(chan func() error, bufferSize)
	chOut := make(chan error, bufferSize)

	return &Group{
		ctx:      ctx,
		wg:       &sync.WaitGroup{},
		limit:    limit,
		duration: duration,
		chIn:     chIn,
		chOut:    chOut,
	}, ctx
}
func (g *Group) RunWorker() {
	g.wg.Add(1)
	go worker(g.ctx, g.wg, g.limit, g.duration, g.chIn, g.chOut)
}

func (g *Group) Go(f func() error) {
	g.chIn <- f
}

func (g *Group) Wait() error {
	close(g.chIn)
	go func() {
		g.wg.Wait()
		close(g.chOut)
	}()

	for res := range g.chOut {
		if res != nil {
			return res
		}
	}

	return nil
}

func worker(ctx context.Context, wg *sync.WaitGroup, limit int, duration time.Duration, jobs <-chan func() error, res chan<- error) {
	ticker := time.NewTicker(duration)

	defer wg.Done()
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for range limit {
				wg.Add(1)
				go func() {
					defer wg.Done()
					job, ok := <-jobs
					if !ok {
						return
					}
					res <- job()
				}()
			}
		}
	}
}
