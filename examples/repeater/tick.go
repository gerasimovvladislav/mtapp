package repeater

import (
	"context"
	"log/slog"
)

var TickID = "Tick"

var Tick = func(ctx context.Context) (cancelFunc context.CancelFunc) {
	ctx, cancelFunc = context.WithCancel(ctx)

	slog.Info("Tick...")

	return
}
