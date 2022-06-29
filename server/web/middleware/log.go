package middleware

import (
	"github.com/gin-gonic/gin"
	log "phanes/collector/logger"
)

func Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.WithFields(log.Fields{
			"request": c.Request.URL.String(),
			"method":  c.Request.Method,
			"ip":      c.ClientIP(),
		}).Info("request")
		// Do something before the request
		c.Next()
		// Do something after the request
	}
}
