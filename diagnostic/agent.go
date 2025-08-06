package diagnostic

import (
	"context"
	"errors"
	"github.com/appellative-ai/core/messaging"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	AgentName = "common:sql:agent/diagnostic/postgres"
)

var (
	agent *agentT
)

func init() {
	NewAgent()
}

type agentT struct {
	running  bool
	poolStat *pgxpool.Stat
	dbClient *pgxpool.Pool
}

func NewAgent() messaging.Agent {
	return newAgent()
}

func newAgent() *agentT {
	a := new(agentT)
	agent = a
	return a
}

// String - identity
func (a *agentT) String() string { return a.Name() }

// Name - agent name
func (a *agentT) Name() string { return AgentName }

// Message - message the agent
func (a *agentT) Message(m *messaging.Message) {
	if m == nil {
		return
	}
	switch m.Name {
	case messaging.ConfigEvent:
		if a.running {
			return
		}
		messaging.UpdateContent[*pgxpool.Pool](&a.dbClient, m)
		return
	case messaging.StartupEvent:
		if a.running {
			return
		}
		a.running = true
		a.run()
		return
	case messaging.ShutdownEvent:
		if !a.running {
			return
		}
		a.running = false
	}
}

// Run - run the agent
func (a *agentT) run() {
}

func (a *agentT) ping(ctx context.Context) error {
	if a.dbClient == nil {
		return errors.New("DbClient is nil")
	}
	return a.dbClient.Ping(ctx)
}

func (a *agentT) stat() error {
	if a.dbClient == nil {
		return errors.New("DbClient is nil")
	}
	a.poolStat = a.dbClient.Stat()
	return nil
}
