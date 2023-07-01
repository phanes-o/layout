package v1

import "github.com/gin-gonic/gin"

// IRouter v1 interface
type IRouter interface {
	Init(r *gin.RouterGroup)
}

var (
	routerList []IRouter
)

// RegisterRouter 注册路由
func RegisterRouter(router IRouter) {
	if router != nil {
		routerList = append(routerList, router)
	}
}

func Init(r *gin.RouterGroup) {
	// 初始化路由
	for _, router := range routerList {
		if router != nil {
			router.Init(r)
		}
	}
}
