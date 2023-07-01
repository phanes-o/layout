package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ResponseOk(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": data,
	})
}
