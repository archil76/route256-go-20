package errgroup

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/net/context"
)

type token struct{ id int }
type Group struct {
	ctx      context.Context
	cancel   context.CancelFunc
	errOnce  sync.Once
	wg       *sync.WaitGroup
	limit    int
	duration time.Duration
	chIn     chan func() error
	chErr    chan error
	chTokens chan token
}

func WithContext(ctx context.Context, limit, bufferSize int, duration time.Duration) (*Group, context.Context) {
	chIn := make(chan func() error, bufferSize)
	chErr := make(chan error, bufferSize)
	chTokens := make(chan token, limit)
	ctx, cancel := context.WithCancel(ctx)

	return &Group{
		ctx:      ctx,
		cancel:   cancel,
		wg:       &sync.WaitGroup{},
		limit:    limit,
		duration: duration,
		chIn:     chIn,
		chErr:    chErr,
		chTokens: chTokens,
	}, ctx
}

func (g *Group) RunWorker() {
	rateCount := 1 + (len(g.chIn) / g.limit)

	g.wg.Add(1)
	go runTicker(g, rateCount)

	go runWorkers(g)
}

func runTicker(g *Group, rateCount int) {
	ticker := time.NewTicker(g.duration)

	defer g.wg.Done()
	defer ticker.Stop()
	defer close(g.chTokens)
	count := 0
	for range rateCount {
		select {
		case <-g.ctx.Done():
			return
		case <-ticker.C:
			for range g.limit {
				count++
				g.chTokens <- token{count}
			}
		}
	}
}

func runWorkers(g *Group) {
	for {
		select {
		case <-g.ctx.Done():
			return
		case token, ok := <-g.chTokens:
			{
				if !ok {
					return
				}
				go worker(token.id, g.ctx, g.wg, g.chIn, g.chErr)
			}
		}
	}
}

func (g *Group) Go(f func() error) {
	g.wg.Add(1)
	g.chIn <- f
}

func (g *Group) Wait() error {
	go func() {
		close(g.chIn)

		g.wg.Wait()

		close(g.chErr)
	}()

	for err := range g.chErr {
		if err != nil {
			g.errOnce.Do(func() {
				g.cancel()
			})

			return err
		}
	}

	return nil
}

func worker(id int, ctx context.Context, wg *sync.WaitGroup, jobs <-chan func() error, errs chan<- error) {
	defer wg.Done()

	fmt.Printf("Job %d start\n", id)

	select {
	case <-ctx.Done():
		fmt.Printf("Job %d was interupted\n", id)
		return
	case job, ok := <-jobs:
		{
			if !ok {
				fmt.Printf("Canal jobs is finished durind the job %d \n", id)
				return
			}

			err := job()
			if err != nil {
				fmt.Printf("Job %d finished with error %e \n", id, err)
				errs <- err
			} else {
				fmt.Printf("Job %d complete\n", id)
			}
		}

	}
}
