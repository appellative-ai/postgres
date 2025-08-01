package retrieval

import (
	"context"
	"errors"
	"github.com/appellative-ai/core/messaging"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"time"
)

const (
	NamespaceName = "sql:postgres:agent/retrieval"
)

var (
	agent *agentT
)

func init() {
	NewAgent()
}

type agentT struct {
	running  bool
	logFunc  func(start time.Time, duration time.Duration, route string, req any, resp any, timeout time.Duration)
	dbClient *pgxpool.Pool
}

func NewAgent() messaging.Agent {
	if agent != nil {
		return agent
	}
	agent = newAgent()
	return agent
}

func newAgent() *agentT {
	a := new(agentT)
	return a
}

// String - identity
func (a *agentT) String() string { return a.Name() }

// Name - agent name
func (a *agentT) Name() string { return NamespaceName }

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
		messaging.UpdateContent[func(start time.Time, duration time.Duration, route string, req any, resp any, timeout time.Duration)](&a.logFunc, m)
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

func (a *agentT) retrieve(ctx context.Context, name, sql string, args []any) (rows pgx.Rows, err error) {
	if a.dbClient == nil {
		return nil, errors.New("DbClient is nil")
	}
	return a.dbClient.Query(ctx, sql, args)
}

func (a *agentT) statusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	return http.StatusInternalServerError
}

func (a *agentT) log(start time.Time, duration time.Duration, route string, req *request, resp *response, ctx context.Context) {
	if a.logFunc == nil {
		return
	}
	var timeout time.Duration
	if d, ok := ctx.Deadline(); ok {
		timeout = time.Until(d)
	}
	a.logFunc(start, duration, route, req, resp, timeout)
}
