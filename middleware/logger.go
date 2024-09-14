package middleware

import (
	"fmt"
	"go-wal/constant"
	"go-wal/pkg/helper"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// get request id
		requestID := c.Request.Header.Get(constant.XRequestId)
		if requestID == "" {
		}
		requestID = helper.GenerateTraceId()
		c.Set(constant.XRequestId, requestID)

		c.Next()
		clientID := c.GetHeader(gin.AuthUserKey)
		end := time.Now()

		status := c.Writer.Status()

		logFields := []interface{}{
			"path", path,
			"client_id", clientID,
			"trace_id", requestID,
			"status", status,
			"method", c.Request.Method,
			"path", path,
			"query", query,
			"ip", c.ClientIP(),
			"profile-agent", c.Request.UserAgent(),
			"time", end.Format(time.RFC3339),
			"latency", fmt.Sprintf("%f", end.Sub(start).Seconds()),
		}

		if status > 400 {
			log.Error().Fields(logFields).Msg("error")
			return
		}
		log.Info().Fields(logFields).Msg("ok")
	}
}
