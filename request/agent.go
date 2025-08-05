package request

import (
	"context"
	"errors"
	"github.com/appellative-ai/core/messaging"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"time"
)

const (
	AgentName              = "common:sql:agent/request/postgres"
	defaultDuration        = time.Second * 3
	StatusTxnBeginError    = int(102) // Transaction processing begin error
	StatusTxnRollbackError = int(103) // Transaction processing rollback error
	StatusTxnCommitError   = int(104) // Transaction processing commit error
	StatusExecError        = int(105) // Execution error, as in a database call
	StatusNotStarted       = int(98)  // Not started
)

var (
	agent    *agentT
	cancelFn = func() {}
)

func init() {
	NewAgent()
}

type agentT struct {
	running  bool
	poolStat *pgxpool.Stat
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

func (a *agentT) exec(ctx context.Context, sql string, args ...any) (tag Result, err error) {
	// Transaction processing.
	if a.dbClient == nil {
		return tag, errors.New("DbClient is nil")
	}
	txn, err0 := a.dbClient.Begin(ctx)
	if err0 != nil {
		return tag, err0
	}
	// Rollback is safe to call even if the tx is already closed, so if
	// the tx commits successfully, this is a no-op
	defer txn.Rollback(ctx)
	cmd, err1 := a.dbClient.Exec(ctx, sql, args)
	if err1 != nil {
		return newResult(cmd), recast(err1)
	}
	err = txn.Commit(ctx)
	return newResult(cmd), err
}

func (a *agentT) statusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	return http.StatusInternalServerError
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
