package notification

import (
	"context"
	"errors"
	"github.com/appellative-ai/core/messaging"
)

// Interface - notification interface
type Interface struct {
	Message        func(ctx context.Context, msg *messaging.Message) error
	ReceiveMessage func(ctx context.Context, name string) (*messaging.Message, error)

	Advise        func(ctx context.Context, msg *messaging.Message) error
	ReceiveAdvice func(ctx context.Context, name string) (*messaging.Message, error)

	Trace func(ctx context.Context, name, task, observation, action string) error
}

// Notifier -
var Notifier = func() *Interface {
	return &Interface{
		Message: func(ctx context.Context, msg *messaging.Message) error {
			return errors.New("not implemented")
		},
		ReceiveMessage: func(ctx context.Context, name string) (*messaging.Message, error) {
			//agent.message(msg)
			return nil, errors.New("not implemented")
		},
		Advise: func(ctx context.Context, msg *messaging.Message) error {
			//agent.advise(msg)
			return errors.New("not implemented")
		},
		ReceiveAdvice: func(ctx context.Context, name string) (*messaging.Message, error) {
			//agent.message(msg)
			return nil, errors.New("not implemented")
		},
		Trace: func(ctx context.Context, name, task, observation, action string) error {
			//agent.trace(name, task, observation, action)
			return nil
		},
	}
}()
