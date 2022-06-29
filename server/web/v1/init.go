package v1

import "github.com/gin-gonic/gin"

// v1 interface

func Init(r *gin.RouterGroup) {
	User.Init(r)
}
