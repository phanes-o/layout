package middleware

import (
	"github.com/gin-gonic/gin"
	log "phanes/collector/logger"
	"phanes/errors"
	"time"
)

func Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		log.WithContext(c.Request.Context()).WithFields(log.Fields{
			"request": c.Request.URL.String(),
			"method":  c.Request.Method,
			"ip":      c.ClientIP(),
		}).Info("request")

		c.Next()

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				errType := errors.GetType(err)
				if errType == errors.None {
					c.JSON(500, gin.H{
						"code":    500,
						"message": "server internal error, more info in logs",
					})
				} else if errType > 3000 && errType < 5000 {
					c.JSON(400, gin.H{
						"code": errType,
						"msg":  err.Error(),
					})
				}
				log.WithContext(c.Request.Context()).WithFields(log.Fields{
					"request":    c.Request.URL.String(),
					"method":     c.Request.Method,
					"ip":         c.ClientIP(),
					"error_type": errType.String(),
					"time_span":  time.Now().Sub(start).String(),
				}).Error(err.Error())
			}
		} else {
			log.WithContext(c.Request.Context()).WithFields(log.Fields{
				"request":   c.Request.URL.String(),
				"method":    c.Request.Method,
				"ip":        c.ClientIP(),
				"time_span": time.Now().Sub(start).String(),
			}).Info("response")
		}
	}
}
