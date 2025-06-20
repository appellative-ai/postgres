package operationstest

import (
	"fmt"
	"github.com/behavioral-ai/core/fmtx"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/postgres/operations"
	"time"
)

// NewService -
func NewService() *operations.Service {
	return &operations.Service{
		Message: func(msg *messaging.Message) bool {
			fmt.Printf("%v  -> %v\n", "message", msg)
			return true
		},
		Advise: func(msg *messaging.Message) *messaging.Status {
			fmt.Printf("%v   -> %v\n", "advise", msg)
			return messaging.StatusOK()
		},
		SubscriptionCreate: func(msg *messaging.Message) {
			fmt.Printf("%v-> %v\n", "subscribe", msg)
		},
		SubscriptionCancel: func(msg *messaging.Message) {
			fmt.Printf("%v   -> %v\n", "cancel", msg)
		},
		Trace: func(name, task, observation, action string) {
			fmt.Printf("%v [%v] [%v] [%v] [%v]", fmtx.FmtRFC3339Millis(time.Now().UTC()), name, task, observation, action)
		},
	}
}
