package repeater

import (
	"context"
	"time"

	"github.com/gerasimovvladislav/mtapp"
)

var ID mtapp.ThreadID = "Repeater"

func Init(interval time.Duration, limit int, run func(ctx context.Context) (cancelFunc context.CancelFunc)) *mtapp.P {
	processor := mtapp.NewProcessor()
	processor.AddThread(mtapp.NewThread(ID, mtapp.NewProcess(run), interval, limit))

	return processor
}
