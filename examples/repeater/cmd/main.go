package main

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/gerasimovvladislav/mtapp"
	"github.com/gerasimovvladislav/mtapp/examples/repeater"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	wg := &sync.WaitGroup{}

	repeat := repeater.Init("Main", time.Millisecond, 3, repeater.Tick)
	mtapp.NewApp(repeat).Start(ctx, wg)

	slog.Info("Application started")

	time.Sleep(300 * time.Millisecond)

	go func() {
		<-ctx.Done()
	}()
	wg.Wait()

	slog.Info("Application has been shutdown gracefully")
}
