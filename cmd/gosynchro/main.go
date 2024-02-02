package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/bmatcuk/doublestar/v3"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"

	"github.com/maxihafer/gosynchro/pkg/logger"
	"github.com/maxihafer/gosynchro/pkg/proxy"
)

var (
	version = "development"
	commit  = ""
	date    = ""
	buildBy = "user"
)

type Config struct {
	remote   string
	port     int
	debug    bool
	json     bool
	patterns cli.StringSlice
}

func (c *Config) Parse() (*proxy.Config, error) {
	_, err := url.Parse(c.remote)
	if err != nil {
		return nil, fmt.Errorf("invalid remote '%s': %w", c.remote, err)
	}

	if c.port < 0 || c.port > 65535 {
		return nil, fmt.Errorf("invalid port '%d'", c.port)
	}

	var files []string
	for _, pattern := range c.patterns.Value() {
		matches, err := doublestar.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid pattern '%s': %w", pattern, err)
		}
		files = append(files, matches...)
	}

	return &proxy.Config{
		Remote: c.remote,
		Port:   c.port,
		Debug:  c.debug,
		Files:  files,
	}, nil
}

func main() {
	var rawConfig Config
	var log zerolog.Logger

	compiled, err := time.Parse(time.RFC3339, date)
	if err != nil {
		panic(err)
	}

	app := &cli.App{
		Name:     "gosynchro",
		Usage:    "A tool to synchronize browser windows",
		Version:  version,
		Compiled: compiled,
		Metadata: map[string]interface{}{
			"commit":  commit,
			"builtBy": buildBy,
		},
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "port",
				Aliases:     []string{"p"},
				Usage:       "Port to listen on",
				Value:       3000,
				Destination: &rawConfig.port,
			},
			&cli.BoolFlag{
				Name:        "verbose",
				Usage:       "Enable verbose logging",
				Value:       false,
				Destination: &rawConfig.debug,
			},
			&cli.BoolFlag{
				Name:        "json",
				Usage:       "Output logs in JSON format",
				Value:       false,
				Destination: &rawConfig.json,
			},
		},
		Before: func(context *cli.Context) error {
			var opts []logger.Option
			if rawConfig.debug {
				opts = append(opts, logger.WithLevel(zerolog.DebugLevel))
			}
			if rawConfig.json {
				opts = append(opts, logger.WithOutput(context.App.Writer))
			} else {
				opts = append(
					opts, logger.WithOutput(
						zerolog.ConsoleWriter{
							Out:        context.App.Writer,
							TimeFormat: time.RFC3339,
						},
					),
				)
			}

			log = logger.New(
				opts...,
			)

			context.Context = log.WithContext(context.Context)

			return nil
		},
		Commands: []*cli.Command{
			&cli.Command{
				Name:  "proxy",
				Usage: "Start a proxy server",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "remote",
						Usage:       "Remote to proxy request to",
						Value:       "http://localhost:8080",
						Destination: &rawConfig.remote,
					},
					&cli.StringSliceFlag{
						Name:        "files",
						Usage:       "Files to watch for changes",
						DefaultText: cli.NewStringSlice("public/**/*.html", "static/*.svg").String(),
						Aliases:     []string{"f"},
						Destination: &rawConfig.patterns,
					},
				},
				Action: func(context *cli.Context) error {
					cfg, err := rawConfig.Parse()
					if err != nil {
						return cli.Exit(err, 1)
					}

					p := &proxy.Proxy{
						Config: cfg,
					}

					if err := p.Start(context.Context); err != nil {
						return cli.Exit(err, 1)
					}

					return nil
				},
			},
			&cli.Command{
				Name:  "reload",
				Usage: "Reload client listening on PORT (default 3000)",
				Action: func(context *cli.Context) error {
					if rawConfig.port < 0 || rawConfig.port > 65535 {
						return cli.Exit("Invalid port", 1)
					}

					resp, err := http.Get("http://localhost:" + strconv.Itoa(rawConfig.port) + "/gosynchro/reload")
					if err != nil {
						panic(err)
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						return cli.Exit("Failed to reload", 1)
					}

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("error")
	}
}
