package operationstest

import (
	"github.com/appellative-ai/core/messaging"
)

func ExampleNewService() {
	s := NewService()
	m := messaging.NewAddressableMessage(messaging.ChannelControl, messaging.ConfigEvent, "core:to/test", "core:from/test")
	s.Message(m)
	s.Advise(m)
	s.SubscriptionCreate(m)
	s.SubscriptionCancel(m)
	s.Trace("core:agent/operations/collective", "task", "going well", "none")

	//Output:
	//message  -> [chan:ctrl] [from:core:from/test] [to:[core:to/test]] [common:core:event/config]
	//advise   -> [chan:ctrl] [from:core:from/test] [to:[core:to/test]] [common:core:event/config]
	//subscribe-> [chan:ctrl] [from:core:from/test] [to:[core:to/test]] [common:core:event/config]
	//cancel   -> [chan:ctrl] [from:core:from/test] [to:[core:to/test]] [common:core:event/config]
	//2025-06-14T22:22:48.229Z [core:agent/operations/collective] [task] [going well] [none]

}
