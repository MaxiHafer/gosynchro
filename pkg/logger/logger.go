package logger

import (
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
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

func StructuredLoggingMiddleware(logger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request = c.Request.WithContext(logger.WithContext(c.Request.Context()))

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		param := gin.LogFormatterParams{}

		param.TimeStamp = time.Now()
		param.Latency = time.Since(start)
		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}
		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		param.BodySize = c.Writer.Size()
		if raw != "" {
			path = path + "?" + raw
		}
		param.Path = path

		var logEvent *zerolog.Event
		if c.Writer.Status() >= 500 {
			logEvent = logger.Error()
		} else {
			logEvent = logger.Debug()
		}

		logEvent.
			Str("client_ip", param.ClientIP).
			Str("method", param.Method).
			Int("status", param.StatusCode).
			Str("path", param.Path).
			Str("latency", param.Latency.String()).
			Msg(param.ErrorMessage)
	}
}
