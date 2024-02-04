package proxy

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/html"

	"github.com/maxihafer/gosynchro/templates"
)

func (p *Proxy) proxyHandler(c *gin.Context) {
	proxy := httputil.NewSingleHostReverseProxy(p.remote)
	proxy.ErrorHandler = p.proxyErrorHandler
	proxy.ModifyResponse = transformHTTPResponse

	proxy.ServeHTTP(c.Writer, c.Request)
}

func (p *Proxy) proxyErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	data := templates.Error{
		Message: err.Error(),
	}

	opErr, ok := err.(*net.OpError)
	if ok && opErr.Op == "dial" {
		data.Title = fmt.Sprintf("Could not connect to %s", p.Config.Remote)
	}

	_ = p.ErrorTemplate.Execute(w, data)
}

func transformHTTPResponse(r *http.Response) error {
	m := r.Header.Get("Content-Type")
	if !strings.HasPrefix(m, "text/html") {
		return nil
	}

	doc, err := html.Parse(r.Body)
	if err != nil {
		return err
	}

	for n := doc.FirstChild; n != nil; n = n.FirstChild {
		if n.Type == html.ElementNode && n.Data == "head" {
			n.AppendChild(
				&html.Node{
					Type: html.ElementNode,
					Data: "script",
					Attr: []html.Attribute{
						{
							Key: "src",
							Val: "/gosynchro/static/gosynchro.js",
						},
					},
				},
			)
			break
		}
	}

	r.Body.Close()
	buf := &bytes.Buffer{}
	r.Body = io.NopCloser(buf)
	if err := html.Render(buf, doc); err != nil {
		return err
	}

	r.Header.Set("Content-Length", fmt.Sprintf("%d", buf.Len()))

	return nil
}
