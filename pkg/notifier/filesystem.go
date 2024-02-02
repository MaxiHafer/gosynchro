package notifier

import (
	"context"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog"

	"github.com/maxihafer/gosynchro/pkg/event"
)

func NewFileSystem(files ...string) (Notifier, error) {
	closer := make(chan struct{})

	fs := &fileSystem{
		watches: files,

		closer: closer,
	}

	return fs, nil
}

var _ Notifier = (*fileSystem)(nil)

type fileSystem struct {
	watches []string

	closer chan struct{}
}

func (f *fileSystem) Notify(ctx context.Context) chan event.Event {
	notify := make(chan event.Event)

	go func() {
		<-f.closer
		close(notify)
	}()

	go f.watch(ctx, notify)

	return notify
}

func (f *fileSystem) watch(ctx context.Context, notify chan event.Event) {
	log := zerolog.Ctx(ctx).With().Str("notifier", "fileSystem").Logger()
	log.Info().Msg("starting notifier")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error().Err(err).Msg("error creating watcher")
		return
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case e, ok := <-watcher.Events:
				if !ok {
					return
				}
				evt := &event.FilesystemEvent{
					Operation: e.Op.String(),
					Path:      e.Name,
				}

				log.Debug().
					Str("operation", evt.Operation).
					Str("path", evt.Path).
					Msg("file system event")

				notify <- evt
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Debug().Err(err).Msg("file system error")
				return
			}
		}
	}()

	for _, watch := range f.watches {
		if err := watcher.Add(watch); err != nil {
			log.Error().Err(err).Str("file", watch).Msg("error registering watcher for file")
			return
		}
	}

	<-f.closer
	log.Info().Msg("closing file system notifier")
}

func (f *fileSystem) Close() {
	close(f.closer)
}
