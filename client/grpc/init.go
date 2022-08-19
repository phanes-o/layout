package grpc

// UserClient example client, in reality, you should replace it with your own grpc client
//var UserClient example.UserService

func Init() func() {
	//UserClient = example.NewUserService("example.UserService", config.MicroService.Client())
	return func() {}
}
