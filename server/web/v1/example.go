package v1

import (
	"github.com/gin-gonic/gin"
	"phanes/bll"
)

type example struct{}

func init() {
	RegisterRouter(&example{})
}

func (e *example) Init(r *gin.RouterGroup) {
	r.GET("/publish", e.publishEvent)
}

func (e *example) publishEvent(c *gin.Context) {
	bll.Example.ExamplePublishEvent(c.Request.Context())
}
