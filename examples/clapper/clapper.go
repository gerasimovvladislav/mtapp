package clapper

import (
	"context"
	"log/slog"
	"time"

	"github.com/gerasimovvladislav/mtapp"
)

func Init(interval time.Duration, limit int) *mtapp.P {
	processor := mtapp.NewProcessor()
	processor.AddThread(mtapp.NewThread("Main", mtapp.NewProcess(func(ctx context.Context) (cancelFunc context.CancelFunc) {
		ctx, cancelFunc = context.WithCancel(ctx)

		slog.Info("Clapp...")

		return
	}), interval, limit))

	return processor
}
