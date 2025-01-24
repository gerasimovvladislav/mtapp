package mtapp

import (
	"context"
	"sync"
	"time"
)

type Processor interface {
	Start(ctx context.Context, wg *sync.WaitGroup)
	Threads() map[ThreadID]*Thread
	Thread(ID ThreadID) *Thread
	AddThread(t *Thread)
	Stop()
}

type P struct {
	mu sync.RWMutex

	threads map[ThreadID]*Thread

	cancelFunc context.CancelFunc
}

func NewProcessor(threads ...*Thread) *P {
	m := make(map[ThreadID]*Thread)

	for _, t := range threads {
		m[t.ID()] = t
	}

	return &P{
		threads: m,
	}
}

func (p *P) AddThread(t *Thread) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.threads[t.ID()] = t
}

func (p *P) Threads() map[ThreadID]*Thread {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.threads
}

func (p *P) Thread(ID ThreadID) *Thread {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.threads[ID]
}

func (p *P) Stop() {
	p.cancelFunc()
}

func (p *P) Start(ctx context.Context, wg *sync.WaitGroup) {
	ctx, p.cancelFunc = context.WithCancel(ctx)

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()

		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				p.mu.RLock()
				for _, t := range p.threads {
					if t.IsPaused() {
						if t.IsRunning() {
							t.Stop()
						}

						continue
					}

					if !t.IsRunning() {
						t.Run(ctx, wg)
					}
				}
				p.mu.RUnlock()
			}
		}
	}(ctx)
}
