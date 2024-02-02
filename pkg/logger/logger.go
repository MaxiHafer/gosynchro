package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type settings struct {
	Level      zerolog.Level
	Output     io.Writer
	TimeFormat string
}

func New(opts ...Option) zerolog.Logger {
	s := &settings{
		Level:  zerolog.InfoLevel,
		Output: os.Stdout,
	}

	for _, opt := range opts {
		opt(s)
	}

	return zerolog.New(s.Output).Level(s.Level).With().Timestamp().Logger()
}

func DefaultWriter() zerolog.ConsoleWriter {
	return zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}
}
