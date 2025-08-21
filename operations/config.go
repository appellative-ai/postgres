package operations

import (
	"github.com/appellative-ai/common/messaging"
)

func (a *agentT) config(m *messaging.Message) {
	if m == nil || m.Name != messaging.ConfigEvent {
		return
	}

}
