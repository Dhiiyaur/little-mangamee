package api

// package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func Logger(log zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// Process the request
		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		Msg := c.Errors.ByType(gin.ErrorTypePrivate).String()

		var logEvent *zerolog.Event
		if c.Writer.Status() >= 500 {
			logEvent = log.Error()
		} else {
			logEvent = log.Info()
		}

		logEvent.
			Str("MESSAGE", Msg).
			Str("METHOD", c.Request.Method).
			Str("PATH", c.Request.URL.Path).
			Str("IP", c.ClientIP()).
			Dur("RESPONSE_TIME", latency).
			Int("STATUS", c.Writer.Status()).
			Send()
	}
}
