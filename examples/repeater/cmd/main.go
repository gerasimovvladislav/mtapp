package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/gerasimovvladislav/mtapp"
	"github.com/gerasimovvladislav/mtapp/examples/repeater"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	repeat := repeater.Init("Main", time.Millisecond, 3, repeater.Tick)
	mtapp.NewApp(repeat).Start(ctx)

	slog.Info("Application has been shutdown gracefully")
}
