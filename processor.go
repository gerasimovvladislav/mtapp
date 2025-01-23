package mtapp

import (
	"context"
	"sync"
	"time"
)

type Processor interface {
	Start(ctx context.Context, wg *sync.WaitGroup)
	AddThread(t *Thread)
	Thread(ID ThreadID) *Thread
	Threads() map[ThreadID]*Thread
	StopThread(ID ThreadID)
	StartThread(ID ThreadID)
	DeleteThread(ID ThreadID)
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

func (p *P) Thread(ID ThreadID) *Thread {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.threads[ID]
}

func (p *P) Threads() map[ThreadID]*Thread {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.threads
}

func (p *P) StartThread(ID ThreadID) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if t := p.threads[ID]; t != nil {
		t.Start()
	}
}

func (p *P) StopThread(ID ThreadID) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if t := p.threads[ID]; t != nil {
		t.Stop()
	}
}

func (p *P) DeleteThread(ID ThreadID) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.threads, ID)
	p.restart()
}

func (p *P) Start(ctx context.Context, wg *sync.WaitGroup) {
	startCtx, cancel := context.WithCancel(ctx)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			p.mu.Lock()
			threads := p.threads
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
							p.DeleteThread(t.ID())
							return
						case <-ticker.C:
							if t.Paused() {
								continue
							}

							t.Work(startCtx)
							if limit > 0 {
								limit--
								if limit == 0 {
									p.DeleteThread(t.ID())

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
				localWg.Wait()
				cancel()
			}
		}
	}()
}
