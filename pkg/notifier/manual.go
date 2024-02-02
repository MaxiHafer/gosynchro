package notifier

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/maxihafer/gosynchro/pkg/event"
)

var _ fmt.Stringer = (*manualEvent)(nil)

type manualEvent struct {
}

func (e manualEvent) String() string {
	return "ManualReloadEvent"
}

func NewManual() *Manual {
	closer := make(chan struct{})
	action := make(chan struct{})

	go func() {
		<-closer
		close(action)
	}()

	return &Manual{
		closer: closer,
		action: action,
	}
}

var _ Notifier = (*Manual)(nil)

type Manual struct {
	action chan struct{}
	closer chan struct{}
}

func (n *Manual) Reload() {
	n.action <- struct{}{}
}

func (n *Manual) Notify(ctx context.Context) chan event.Event {
	notify := make(chan event.Event)

	go func() {
		<-n.closer
		close(notify)
	}()

	go n.watch(ctx, notify)

	return notify
}

func (n *Manual) watch(ctx context.Context, notify chan event.Event) {
	log := zerolog.Ctx(ctx).With().Str("notifier", "manual").Logger()
	log.Info().Msg("starting notifier")

	go func() {
		for {
			select {
			case <-n.action:
				notify <- &event.ReloadEvent{}
				log.Debug().Msg("received event")
			}
		}
	}()

	<-n.closer
	log.Info().Msg("closing notifier")
}

func (m *Manual) Close() {
	close(m.closer)
}
