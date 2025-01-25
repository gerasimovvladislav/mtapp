package mtapp

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

type Thread struct {
	mu sync.RWMutex

	id        ThreadID
	process   *Process
	interval  time.Duration
	limit     int
	limitable bool
	paused    bool

	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewThread(ID ThreadID, p *Process, paused bool, interval time.Duration, limit int) *Thread {
	return &Thread{
		id:        ID,
		process:   p,
		paused:    paused,
		limitable: limit > 0,
		interval:  interval,
		limit:     limit,
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

func (t *Thread) IsPaused() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.paused
}

func (t *Thread) IsRunning() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.cancelFunc != nil
}

func (t *Thread) Pause() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.paused = true
}

func (t *Thread) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.Pause()
	if t.cancelFunc != nil {
		t.cancelFunc()
		slog.Debug("Thread stopped", slog.String("id", string(t.id)))
		t.cancelFunc = nil
	}
}

func (t *Thread) Start(ctx context.Context) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.ctx = ctx
	t.paused = false
}

func (t *Thread) Run(ctx context.Context, wg *sync.WaitGroup) {
	ctx, t.cancelFunc = context.WithCancel(ctx)

	t.Start(ctx)

	wg.Add(1)
	defer wg.Done()

	ticker := time.NewTicker(t.interval)
	defer ticker.Stop()

	slog.Debug("Thread started", slog.String("id", string(t.id)))

	for {
		select {
		case <-t.ctx.Done():
			return
		case <-ticker.C:
			if t.IsPaused() {
				continue
			}

			t.mu.Lock()
			if t.limitable {
				if t.limit > 0 {
					t.work(t.ctx)
					t.limit--
					if t.limit == 0 {
						return
					}
				}
			} else {
				t.work(t.ctx)
			}
			t.mu.Unlock()
		}
	}
}

func (t *Thread) work(ctx context.Context) {
	if !t.process.IsRunning() {
		start := time.Now()
		t.process.Run(ctx)
		recordProcessLatency(t.id, time.Since(start))
	}
}
