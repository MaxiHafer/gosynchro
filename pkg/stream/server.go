package stream

import (
	"io"

	"github.com/gin-gonic/gin"

	"github.com/maxihafer/gosynchro/pkg/event"
)

const ContextKey = "gosynchro-client"

func NewServer(events chan event.Event) *Server {
	srv := &Server{
		Events:        events,
		NewClients:    make(chan chan event.Event),
		ClosedClients: make(chan chan event.Event),
		Clients:       make(map[chan event.Event]struct{}),
	}

	go srv.listen()

	return srv
}

type Server struct {
	Events        chan event.Event
	NewClients    chan chan event.Event
	ClosedClients chan chan event.Event
	Clients       map[chan event.Event]struct{}
}

func (s *Server) listen() {
	for {
		select {
		case client := <-s.NewClients:
			s.Clients[client] = struct{}{}
		case client := <-s.ClosedClients:
			delete(s.Clients, client)
			close(client)
		case evt := <-s.Events:
			for client := range s.Clients {
				client <- evt
			}
		}
	}
}

func ClientStreamFromContext(c *gin.Context) (chan string, bool) {
	v, ok := c.Get(ContextKey)
	if !ok {
		return nil, false
	}
	clientChan, ok := v.(chan string)
	if !ok {
		return nil, false
	}
	return clientChan, true
}

func (s *Server) StreamEvents(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	client := make(chan event.Event)
	s.NewClients <- client
	defer func() {
		s.ClosedClients <- client
	}()

	c.Stream(
		func(w io.Writer) bool {
			if msg, ok := <-client; ok {
				writeEvent(c, msg)
				return true
			}
			return false
		},
	)
}

func writeEvent(c *gin.Context, event event.Event) {
	c.SSEvent(event.Name(), event)
}
