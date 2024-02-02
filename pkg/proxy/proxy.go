package proxy

import (
	"fmt"
	"net"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"

	"github.com/maxihafer/gosynchro"
	"github.com/maxihafer/gosynchro/pkg/notifier"
	"github.com/maxihafer/gosynchro/pkg/stream"
)

type Proxy struct {
	*Config
	manualNotifier *notifier.Manual
}

func (p *Proxy) Start(ctx context.Context) error {
	if p.Config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	remote, err := url.Parse(p.Config.Remote)
	if err != nil {
		return err
	}

	p.manualNotifier = notifier.NewManual()

	aggregateNotifier := notifier.NewAggregate(
		p.manualNotifier,
	)

	if len(p.Config.Files) > 0 {
		fs, err := notifier.NewFileSystem(p.Config.Files...)
		if err != nil {
			return err
		}
		aggregateNotifier.Add(fs)
	}

	stream := stream.NewServer(aggregateNotifier.Notify(ctx))

	engine := gin.Default()
	engine.SetTrustedProxies(nil)

	engine.GET("/gosynchro", stream.StreamEvents())
	engine.StaticFS("/gosynchro/static", http.FS(gosynchro.StaticFS))
	engine.GET(
		"/gosynchro/reload", func(c *gin.Context) {
			p.manualNotifier.Reload()
			c.JSON(http.StatusOK, gin.H{"message": "Reloaded"})
		},
	)

	engine.NoRoute(p.proxyHandler(remote))

	addr := net.JoinHostPort("", fmt.Sprintf("%d", p.Port))
	if err := engine.Run(addr); err != nil {
		return err
	}

	return nil
}
