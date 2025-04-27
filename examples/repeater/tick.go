package repeater

import (
	"context"
	"fmt"
)

var TickID = "Tick"

var Tick = func(ctx context.Context) (cancelFunc context.CancelFunc) {
	ctx, cancelFunc = context.WithCancel(ctx)

	fmt.Println("Tick...")

	return
}
