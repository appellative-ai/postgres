package notification

import (
	"github.com/appellative-ai/core/messaging"
	"github.com/appellative-ai/core/std"
)

// Interface - notification interface
type Interface struct {
	Message        func(msg *messaging.Message) bool
	ReceiveMessage func(name string) *messaging.Message

	Advise        func(msg *messaging.Message) *std.Status
	ReceiveAdvice func(name string) *messaging.Message

	Trace func(name, task, observation, action string)
}

// Notifier -
var Notifier = func() *Interface {
	return &Interface{
		Message: func(msg *messaging.Message) bool {
			//agent.message(msg)
			return true
		},
		ReceiveMessage: func(name string) *messaging.Message {
			//agent.message(msg)
			return nil
		},
		Advise: func(msg *messaging.Message) *std.Status {
			//agent.advise(msg)
			return std.StatusOK
		},
		ReceiveAdvice: func(name string) *messaging.Message {
			//agent.message(msg)
			return nil
		},
		Trace: func(name, task, observation, action string) {
			//agent.trace(name, task, observation, action)
		},
	}
}()
