package mtapp

import (
	"context"
	"sync"
	"time"
)

type Processor interface {
	Start(ctx context.Context, wg *sync.WaitGroup)
	AddThread(t *Thread)
	DelThread(t *Thread)
}

type P struct {
	mu      sync.RWMutex
	threads map[ThreadID]*Thread
	buf     chan struct{}
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

func NewProcessor() *P {
	return &P{
		threads: make(map[ThreadID]*Thread),
		buf:     make(chan struct{}, 1),
	}
}

func (p *P) restart() {
	select {
	case p.buf <- struct{}{}:
	default:
	}
}

func (p *P) AddThread(t *Thread) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.threads[t.ID()] = t
	p.restart()
}

func (p *P) DelThread(t *Thread) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.threads, t.ID())
	p.restart()
}

func (p *P) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			p.mu.Lock()
			threads := make(map[ThreadID]*Thread)
			for id, t := range p.threads {
				threads[id] = t
			}
			p.mu.Unlock()

			startCtx, cancel := context.WithCancel(ctx)
			p.mu.Lock()
			p.cancel = cancel
			p.mu.Unlock()

			localWg := &sync.WaitGroup{}
			for _, t := range threads {
				limit := t.Limit()
				localWg.Add(1)
				go func(t *Thread) {
					defer localWg.Done()

					ticker := time.NewTicker(t.Interval())
					defer ticker.Stop()

					for {
						select {
						case <-startCtx.Done():
							return
						case <-ticker.C:
							t.Work(startCtx)
							if limit > 0 {
								limit--
								if limit == 0 {
									return
								}
							}
						}
					}
				}(t)
			}

			select {
			case <-ctx.Done():
				cancel()
				localWg.Wait()
				return
			case <-p.buf:
				cancel()
				localWg.Wait()
			}
		}
	}()
}
