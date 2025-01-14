package main

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/gerasimovvladislav/mtapp"
	"github.com/gerasimovvladislav/mtapp/examples/clapper"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	wg := &sync.WaitGroup{}

	mtapp.NewApp(clapper.Init(time.Second, 3)).Start(ctx, wg)

	slog.Info("Application started")

	go func() {
		<-ctx.Done()
	}()
	wg.Wait()

	slog.Info("Application has been shutdown gracefully")
}
