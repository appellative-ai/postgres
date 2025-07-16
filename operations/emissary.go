package operations

import (
	"github.com/appellative-ai/core/messaging"
)

// emissary attention
func emissaryAttend(a *agentT) {
	var paused = false
	if paused {
	}
	//a.ticker.Start(-1)
	for {
		select {
		case <-a.ticker.T.C:
			// TODO: retrieval collective for new messages and advice
		default:
		}
		select {
		case msg := <-a.emissary.C:
			switch msg.Name {
			case messaging.PauseEvent:
				paused = true
			case messaging.ResumeEvent:
				paused = false
			case messaging.ShutdownEvent:
				a.emissaryFinalize()
				return
			default:
			}
		default:
		}
	}
}
