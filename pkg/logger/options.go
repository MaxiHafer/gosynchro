package logger

import (
	"io"

	"github.com/rs/zerolog"
)

type Option func(settings *settings)

func WithOutput(writer io.Writer) Option {
	return func(settings *settings) {
		settings.Output = writer
	}
}

func WithLevel(level zerolog.Level) Option {
	return func(settings *settings) {
		settings.Level = level
	}
}
