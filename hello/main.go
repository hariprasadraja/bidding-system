package main

import (
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro"
	"hello/handler"
	"hello/subscriber"

	hello "hello/proto/hello"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("com.foo.srv.hello"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	hello.RegisterHelloHandler(service.Server(), new(handler.Hello))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("com.foo.srv.hello", service.Server(), new(subscriber.Hello))

	// Register Function as Subscriber
	micro.RegisterSubscriber("com.foo.srv.hello", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
