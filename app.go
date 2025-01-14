package mtapp

import (
	"context"
	"sync"
)

type App struct {
	processors []Processor
}

func NewApp(processors ...Processor) *App {
	return &App{
		processors: processors,
	}
}

func (a *App) Start(ctx context.Context, wg *sync.WaitGroup) {
	for _, p := range a.processors {
		p.Start(ctx, wg)
	}
}
