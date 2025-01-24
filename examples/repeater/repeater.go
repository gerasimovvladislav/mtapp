package repeater

import (
	"context"
	"time"

	"github.com/gerasimovvladislav/mtapp"
)

func Init(tid mtapp.ThreadID, interval time.Duration, limit int, run func(ctx context.Context) (cancelFunc context.CancelFunc)) *mtapp.P {
	processor := mtapp.NewProcessor(
		mtapp.NewThread(
			tid,
			mtapp.NewProcess(run),
			false,
			interval,
			limit,
		),
	)

	return processor
}
