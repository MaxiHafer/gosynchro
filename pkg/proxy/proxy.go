package proxy

import (
	"fmt"
	"net"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"golang.org/x/net/context"

	"github.com/maxihafer/gosynchro/pkg/logger"
	"github.com/maxihafer/gosynchro/pkg/notifier"
	"github.com/maxihafer/gosynchro/pkg/stream"
)

func NewFromConfig(c *Config) (*Proxy, error) {
	p := &Proxy{Config: c}

	remote, err := url.Parse(c.Remote)
	if err != nil {
		return nil, fmt.Errorf("invalid remote '%s': %w", c.Remote, err)
	}
	p.remote = remote
	p.listenAddr = net.JoinHostPort("", fmt.Sprintf("%d", c.Port))

	p.manualNotifier = notifier.NewManual()

	p.aggregateNotifier = notifier.NewAggregate(p.manualNotifier)

	if len(c.Files) > 0 {
		fs, err := notifier.NewFileSystem(c.Files...)
		if err != nil {
			return nil, err
		}
		p.aggregateNotifier.Add(fs)
	}

	return p, nil
}

type Proxy struct {
	*Config

	aggregateNotifier *notifier.Aggregate
	manualNotifier    *notifier.Manual
	remote            *url.URL
	listenAddr        string

	streamServer *stream.Server
}

func (p *Proxy) Start(ctx context.Context) error {
	gin.SetMode(gin.ReleaseMode)

	log := zerolog.Ctx(ctx)

	stream := stream.NewServer(p.aggregateNotifier.Notify(ctx))

	engine := gin.New()
	engine.Use(logger.StructuredLoggingMiddleware(log))
	engine.Use(gin.Recovery())
	engine.SetTrustedProxies(nil)

	engine.GET("/gosynchro", stream.StreamEvents)
	engine.StaticFS("/gosynchro/static", http.FS(p.StaticFS))
	engine.GET(
		"/gosynchro/reload", func(c *gin.Context) {
			p.manualNotifier.Reload()
			c.JSON(http.StatusOK, gin.H{"message": "Reloaded"})
		},
	)
	engine.NoRoute(p.proxyHandler)

	if err := engine.Run(p.listenAddr); err != nil {
		return err
	}

	return nil
}
