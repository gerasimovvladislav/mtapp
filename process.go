package mtapp

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

type Process struct {
	mu sync.RWMutex

	id ProcessID

	run        func(ctx context.Context) context.CancelFunc
	cancelFunc context.CancelFunc

	isRunning bool
}

func NewProcess(run func(ctx context.Context) context.CancelFunc) *Process {
	return &Process{
		id:        ProcessID(uuid.New()),
		run:       run,
		isRunning: false,
	}
}

func (p *Process) ID() ProcessID {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.id
}

func (p *Process) IsRunning() bool {
	return p.isRunning
}

func (p *Process) Run(ctx context.Context) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.IsRunning() {
		return
	}

	p.isRunning = true
	if p.run != nil {
		p.cancelFunc = p.run(ctx)
	}
	p.isRunning = false
}

func (p *Process) Kill() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.isRunning {
		p.cancelFunc()
		p.isRunning = false
	}
}
