package subscriber

import (
	"context"
	"github.com/micro/go-micro/util/log"

	hello "hello/proto/hello"
)

type Hello struct{}

func (e *Hello) Handle(ctx context.Context, msg *hello.Message) error {
	log.Log("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *hello.Message) error {
	log.Log("Function Received message: ", msg.Say)
	return nil
}
