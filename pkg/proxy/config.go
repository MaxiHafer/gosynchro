package proxy

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/url"
)

type Config struct {
	Remote string
	Port   int
	Debug  bool
	Files  []string

	StaticFS      fs.FS
	ErrorTemplate *template.Template
}

func (c *Config) Validate() error {
	_, err := url.Parse(c.Remote)
	if err != nil {
		return fmt.Errorf("invalid remote '%s': %w", c.Remote, err)
	}

	if c.Port < 0 || c.Port > 65535 {
		return fmt.Errorf("invalid port '%d'", c.Port)
	}

	return nil
}
