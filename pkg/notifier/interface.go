package notifier

import (
	"context"

	"github.com/maxihafer/gosynchro/pkg/event"
)

type Notifier interface {
	Notify(ctx context.Context) chan event.Event
	Close()
}
