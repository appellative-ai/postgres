package operations

import (
	"fmt"
	"github.com/appellative-ai/core/messaging"
	"github.com/appellative-ai/postgres/request"
	"github.com/appellative-ai/postgres/retrieval"
	"sync/atomic"
	"time"
)

const (
	AgentName = "common:sql:agent/operations/postgres"
	duration  = time.Second * 30
)

var (
	agent *agentT
	Agent messaging.Agent
)

func init() {
	agent = newAgent()
	Agent = agent
}

type agentT struct {
	running atomic.Bool
	state   *operationsT
	agents  *messaging.Exchange

	ticker   *messaging.Ticker
	emissary *messaging.Channel
}

func newAgent() *agentT {
	a := new(agentT)
	a.running.Store(false)
	a.agents = messaging.NewExchange()
	a.agents.Register(request.NewAgent())
	a.agents.Register(retrieval.NewAgent())
	agent = a

	a.ticker = messaging.NewTicker(messaging.ChannelEmissary, duration)
	a.emissary = messaging.NewEmissaryChannel()
	a.configureAgents()
	return a
}

// String - identity
func (a *agentT) String() string { return a.Name() }

// Name - agent identifier
func (a *agentT) Name() string { return AgentName }

// Message - message the agent
func (a *agentT) Message(m *messaging.Message) {
	if m == nil {
		return
	}
	switch m.Name {
	case messaging.ConfigEvent:
		a.config(m)
		return
	case messaging.StartupEvent:
		if a.running.Load() {
			return
		}
		a.running.Store(true)
		a.run()
		return
	case messaging.ShutdownEvent:
		if !a.running.Load() {
			return
		}
		a.running.Store(false)
	case messaging.PauseEvent:
		//a.enabled.Store(false)
		//a.events.empty()
	case messaging.ResumeEvent:
		//a.enabled.Store(true)
	}
	switch m.Channel() {
	case messaging.ChannelControl, messaging.ChannelEmissary:
		a.emissary.C <- m
	default:
		fmt.Printf("limiter - invalid channel %v\n", m)
	}
}

// Run - run the agent
func (a *agentT) run() {
	go emissaryAttend(a)
}

func (a *agentT) emissaryFinalize() {
	a.emissary.Close()
	a.ticker.Stop()
}

func (a *agentT) configureAgents() {
	///a.agents.Broadcast(private.NewInterfaceMessage(&private.Interface{
	//	Representation: representation,
	//	Context:        context,
	//	Thing:          thing,
	//	Relation:       relation,
	//}))
}
