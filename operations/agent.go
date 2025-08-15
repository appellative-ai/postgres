package operations

import (
	"context"
	"errors"
	"fmt"
	"github.com/appellative-ai/core/messaging"
	"github.com/appellative-ai/postgres/request"
	"github.com/appellative-ai/postgres/retrieval"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"sync/atomic"
	"time"
)

const (
	AgentName       = "common:sql:agent/operations/postgres"
	duration        = time.Second * 30
	timeoutDuration = time.Millisecond * 1500
)

var (
	agent *agentT
	Agent messaging.Agent
)

func init() {
	Agent = newAgent()
}

type agentT struct {
	running atomic.Bool
	agents  *messaging.Exchange

	poolStat *pgxpool.Stat
	logFunc  func(start time.Time, duration time.Duration, route string, req any, resp any, timeout time.Duration)
	dbClient *pgxpool.Pool

	ticker   *messaging.Ticker
	emissary *messaging.Channel
}

func newAgent() *agentT {
	a := new(agentT)
	agent = a
	a.running.Store(false)
	a.agents = messaging.NewExchange()
	a.agents.Register(request.NewAgent())
	a.agents.Register(retrieval.NewAgent())

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

func (a *agentT) clientStartup(cfg map[string]string) error {
	err := clientStartup(cfg)
	if err != nil {
		return err
	}
	m := messaging.NewConfigMessage(a.dbClient)
	a.agents.Broadcast(m)
	return nil
}

func (a *agentT) configLogging(log func(start time.Time, duration time.Duration, route string, req any, resp any, timeout time.Duration)) {
	if log == nil {
		return
	}
	a.logFunc = log
	m := messaging.NewConfigMessage(log)
	a.agents.Broadcast(m)
}

func (a *agentT) ping(ctx context.Context) error {
	if a.dbClient == nil {
		return errors.New("DbClient is nil")
	}
	start := time.Now().UTC()
	err := a.dbClient.Ping(ctx)
	a.log(start, time.Since(start), pingRouteName, newRequest("", ""), newResponse(a.statusCode(err)), ctx)
	return err
}

func (a *agentT) stat() error {
	if a.dbClient == nil {
		return errors.New("DbClient is nil")
	}
	a.poolStat = a.dbClient.Stat()
	return nil
}

func (a *agentT) statusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	return http.StatusInternalServerError
}

func (a *agentT) log(start time.Time, duration time.Duration, route string, req *requestT, resp *response, ctx context.Context) {
	if a.logFunc == nil {
		return
	}
	var timeout time.Duration
	if d, ok := ctx.Deadline(); ok {
		timeout = time.Until(d)
	}
	a.logFunc(start, duration, route, req, resp, timeout)
}
