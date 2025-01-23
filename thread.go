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
	paused   bool
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

func (t *Thread) Paused() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.paused
}

func (t *Thread) Start() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.paused = false
}

func (t *Thread) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.paused = true
}

func (t *Thread) Work(ctx context.Context) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	t.paused = false

	if !t.process.IsRunning() {
		start := time.Now()
		t.process.Run(ctx)
		recordProcessLatency(t.id, time.Since(start))
	}
}
