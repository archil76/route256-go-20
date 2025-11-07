package ratelimiter

import (
	"context"
	"errors"
	"time"
)

type token struct{}

type RateLimiter struct {
	ctx      context.Context
	cancel   context.CancelFunc
	chTokens chan token
	limit    int
	duration time.Duration
}

func WithContext(ctx context.Context, limit int, duration time.Duration) (*RateLimiter, error) {
	if limit <= 0 {
		return nil, errors.New("limit must be greater than zero")
	}

	chTokens := make(chan token, limit)
	ctx, cancel := context.WithCancel(ctx)

	return &RateLimiter{
		ctx:      ctx,
		cancel:   cancel,
		chTokens: chTokens,
		limit:    limit,
		duration: duration,
	}, nil
}

func (r *RateLimiter) Start() {
	for range r.limit {
		r.chTokens <- token{}
	}
	ticker := time.NewTicker(r.duration)

	go func() {
		defer ticker.Stop()
		defer close(r.chTokens)
		for {
			select {
			case <-r.ctx.Done():
				return
			case <-ticker.C:
				for range r.limit {
					select {
					case r.chTokens <- token{}:
					default:
					}

				}

			}
		}
	}()
}

func (r *RateLimiter) Wait() {
	select {
	case <-r.ctx.Done():
		return
	case <-r.chTokens:
	}
}

func (r *RateLimiter) Stop() {
	r.cancel()
}
