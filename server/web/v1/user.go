package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/phanes-o/proto/dto"
	"github.com/phanes-o/proto/primitive"
	"phanes/bll"
)

var User = &user{}

type user struct{}

func (a *user) Init(r *gin.RouterGroup) {
	u := r.Group("user")
	{
		u.POST("register", a.register)
		u.DELETE("/:value", a.delete)
	}
}

func (a *user) register(c *gin.Context) {
	var user = &dto.CreateUserRequest{}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.Error(err)
		return
	}

	if err := bll.User.Create(c.Request.Context(), user); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
	})
}

func (a *user) delete(c *gin.Context) {
	d := &primitive.Int64{}

	if err := c.ShouldBindJSON(&d); err != nil {
		c.Error(err)
		return
	}

	if err := bll.User.Delete(c.Request.Context(), d); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
	})
}
