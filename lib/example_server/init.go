package example_server

import (
	"go-micro.dev/v4/server"
	"log"
	"net"
)

func NewServer() server.Server {
	return exampleServer{}
}

type exampleServer struct {
	opts *server.Options
}

func (e exampleServer) Init(option ...server.Option) error {
	for _, opt := range option {
		opt(e.opts)
	}
	// Do some other init here
	return nil
}

func (e exampleServer) handleConn(conn net.Conn) {
	for {
		buff := make([]byte, 1024)
		_, err := conn.Read(buff)
		if err != nil {
			log.Println(err)
		}
		// do something
	}
}

func (e exampleServer) Options() server.Options {
	//TODO implement me
	panic("implement me")
}

func (e exampleServer) Handle(handler server.Handler) error {
	//TODO implement me
	panic("implement me")
}

func (e exampleServer) NewHandler(i interface{}, option ...server.HandlerOption) server.Handler {
	//TODO implement me
	panic("implement me")
}

func (e exampleServer) NewSubscriber(s string, i interface{}, option ...server.SubscriberOption) server.Subscriber {
	//TODO implement me
	panic("implement me")
}

func (e exampleServer) Subscribe(subscriber server.Subscriber) error {
	//TODO implement me
	panic("implement me")
}

func (e exampleServer) Start() error {
	listen, err := net.Listen("tcp", ":9090")
	if err != nil {
		return err
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println(err)
		}
		e.handleConn(conn)
	}
}

func (e exampleServer) Stop() error {
	//TODO implement me
	panic("implement me")
}

func (e exampleServer) String() string {
	//TODO implement me
	panic("implement me")
}
