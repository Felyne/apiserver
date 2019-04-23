package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for incoming header, use it if exists
		requestId := c.Request.Header.Get("X-Request-Id")
		if requestId == "" {
			guid := xid.New()
			requestId = guid.String()
		}
		c.Set("X-Request-Id", requestId)
		c.Writer.Header().Set("X-Request-Id", requestId)
		c.Next()
	}
}
