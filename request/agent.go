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
	NamespaceName          = "sql:postgres:agent/request"
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
	logFunc  func(start time.Time, duration time.Duration, req any, resp any, timeout time.Duration)
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
		messaging.UpdateContent[func(start time.Time, duration time.Duration, req any, resp any, timeout time.Duration)](&a.logFunc, m)
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

func (a *agentT) exec(ctx context.Context, sql string, args ...any) (tag Response, status *messaging.Status) {
	// Transaction processing.
	if a.dbClient == nil {
		return tag, messaging.NewStatus(messaging.StatusInvalidArgument, errors.New("DbClient is nil"))
	}
	txn, err0 := a.dbClient.Begin(ctx)
	if err0 != nil {
		return tag, messaging.NewStatus(StatusTxnBeginError, err0)
	}
	// Rollback is safe to call even if the tx is already closed, so if
	// the tx commits successfully, this is a no-op
	defer txn.Rollback(ctx)
	cmd, err := a.dbClient.Exec(ctx, sql, args)
	if err != nil {
		return newResponse(cmd), messaging.NewStatus(messaging.StatusInvalidArgument, recast(err))
	}
	err = txn.Commit(ctx)
	if err != nil {
		status = messaging.NewStatus(StatusTxnCommitError, err)
	} else {
		status = messaging.StatusOK()
	}
	return newResponse(cmd), status
}

func (a *agentT) statusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	return http.StatusInternalServerError
}

/*
	func (a *agentT) setTimeout(ctx context.Context) (context.Context, func()) {
		if ctx == nil {
			ctx = context.Background()
		}
		if d, ok := ctx.Deadline(); ok {
			a.state.Until = time.Until(d)
			return ctx, cancelFn
		}
		if a.state.Timeout <= 0 {
			return ctx, cancelFn
		}
		return context.WithTimeout(ctx, a.state.Timeout)

}
*/
func (a *agentT) log(start time.Time, duration time.Duration, req *request, resp *logResponse, ctx context.Context) {
	if a.logFunc == nil {
		return
	}
	var timeout time.Duration
	if d, ok := ctx.Deadline(); ok {
		timeout = time.Until(d)
	}
	a.logFunc(start, duration, req, resp, timeout)
}
