package errgroup

import (
	"sync"
	"time"

	"golang.org/x/net/context"
)

type Group struct {
	wg    *sync.WaitGroup
	limit int
	chIn  chan func() error
	chOut chan error
}

func WithContext(ctx context.Context, limit, bufferSize int) (*Group, context.Context) {
	chIn := make(chan func() error, bufferSize)
	chOut := make(chan error, bufferSize)

	return &Group{
		wg:    &sync.WaitGroup{},
		limit: limit,
		chIn:  chIn,
		chOut: chOut,
	}, ctx
}
func (g *Group) RunWorker() {

	for i := 0; i < g.limit; i++ {
		g.wg.Add(1)
		go worker(g.wg, g.chIn, g.chOut)
	}
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

func worker(wg *sync.WaitGroup, jobs <-chan func() error, res chan<- error) {
	defer wg.Done()
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}
			res <- job()
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}
