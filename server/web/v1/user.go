package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/phanes-o/proto/base"
	"github.com/phanes-o/proto/dto"
	"phanes/bll"
	"phanes/errors"
)

var User = &user{}

func init() {
	RegisterRouter(&user{})
}

type user struct{}

func (a *user) Init(r *gin.RouterGroup) {
	u := r.Group("user")
	{
		u.POST("register", a.register)
		u.DELETE("/:value", a.delete)
	}
}

func (a *user) register(c *gin.Context) {
	var u = &dto.CreateUserRequest{}

	if err := c.ShouldBindJSON(&u); err != nil {
		c.Error(errors.ErrParamsParse.Warp(err, "register unmarshal json error"))
		return
	}

	if u.Username == "" {
		c.Error(errors.BadRequest.New("username is required"))
		return
	}

	if u.Password == "" {
		c.Error(errors.BadRequest.New("password is required"))
		return
	}
	if err := bll.User.Create(c.Request.Context(), u); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
	})
}

func (a *user) delete(c *gin.Context) {
	d := &base.Int64{}

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
