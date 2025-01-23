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

	repeat.AddThread(mtapp.NewThread("Other", mtapp.NewProcess(func(ctx context.Context) (cancelFunc context.CancelFunc) {
		ctx, cancelFunc = context.WithCancel(ctx)

		slog.Info("New Tick...")

		return
	}), time.Second, 3))

	go func() {
		<-ctx.Done()
	}()
	wg.Wait()

	slog.Info("Application has been shutdown gracefully")
}
