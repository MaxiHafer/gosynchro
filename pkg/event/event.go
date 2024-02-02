package event

const EventNamePrefix = "gosynchro/"

type Event interface {
	Name() string
}

var _ Event = (*FilesystemEvent)(nil)

type FilesystemEvent struct {
	Operation string `json:"operation"`
	Path      string `json:"path"`
}

func (f FilesystemEvent) Name() string {
	return EventNamePrefix + "filesystem"
}

var _ Event = (*ReloadEvent)(nil)

type ReloadEvent struct{}

func (r ReloadEvent) Name() string {
	return EventNamePrefix + "reload"
}

var _ Event = (*ConnectEvent)(nil)

type ConnectEvent struct {
	Remote string `json:"remote"`
}

func (c ConnectEvent) Name() string {
	return EventNamePrefix + "connect"
}
