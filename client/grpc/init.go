package grpc

import (
	"github.com/phanes-o/proto/example"
	"phanes/server"
)

// UserClient example client, in reality, you should replace it with your own grpc client
var UserClient example.UserService

func Init() func() {
	UserClient = example.NewUserService("example.UserService", server.Service.Client())
	return func() {}
}
