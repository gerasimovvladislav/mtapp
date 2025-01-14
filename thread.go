package mtapp

import (
	"context"
	"sync"
	"time"
)

type Thread struct {
	mu       sync.RWMutex
	id       ThreadID
	process  *Process
	interval time.Duration
	limit    int
}

func NewThread(ID ThreadID, p *Process, interval time.Duration, limit int) *Thread {
	return &Thread{
		id:       ID,
		process:  p,
		interval: interval,
		limit:    limit,
	}
}

func (t *Thread) Interval() time.Duration {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.interval
}

func (t *Thread) Limit() int {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.limit
}

func (t *Thread) ID() ThreadID {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.id
}

func (t *Thread) Work(ctx context.Context) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if !t.process.IsRunning() {
		t.process.Run(ctx)
	}
}
