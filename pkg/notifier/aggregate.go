package notifier

import (
	"context"

	"github.com/maxihafer/gosynchro/pkg/event"
)

func NewAggregate(notifiers ...Notifier) *Aggregate {
	closer := make(chan struct{})

	return &Aggregate{
		notifiers: notifiers,
		closer:    closer,
	}
}

var _ Notifier = (*Aggregate)(nil)

type Aggregate struct {
	notifiers []Notifier
	closer    chan struct{}
}

func (n *Aggregate) Add(notifier Notifier) {
	n.notifiers = append(n.notifiers, notifier)
}

func (n *Aggregate) Notify(ctx context.Context) chan event.Event {
	notify := make(chan event.Event)

	go func() {
		<-n.closer
		close(notify)
	}()

	for _, notifier := range n.notifiers {
		curN := notifier
		go func() {
			for msg := range curN.Notify(ctx) {
				notify <- msg
			}
		}()
	}

	return notify
}

func (n *Aggregate) Close() {
	for _, notifier := range n.notifiers {
		notifier.Close()
	}

	close(n.closer)
}
