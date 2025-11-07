package errgroup

import (
	"context"
	"sync"
)

type Group struct {
	ctx     context.Context
	cancel  context.CancelFunc
	errOnce sync.Once
	wg      *sync.WaitGroup
	err     error
}

func WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)

	return &Group{
		ctx:    ctx,
		cancel: cancel,
		wg:     &sync.WaitGroup{},
	}, ctx
}

func (g *Group) Go(f func() error) {
	if g.err != nil {
		return
	}
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()
		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.cancel()
				g.err = err
			})
		}
	}()
}

func (g *Group) Wait() error {
	g.wg.Wait()
	return g.err
}
