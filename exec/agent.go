package exec

import (
	"context"
	"errors"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/postgres/private"
	"net/http"
	"time"
)

const (
	NamespaceName          = "sql:postgres:agent/exec"
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
	running bool
	state   *private.Configuration
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
	a.state = private.NewConfiguration(defaultDuration)
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
	if !a.running {
		if m.Name == messaging.ConfigEvent {
			a.configure(m)
			return
		}
		if m.Name == messaging.StartupEvent {
			a.run()
			a.running = true
			return
		}
		return
	}
	if m.Name == messaging.ShutdownEvent {
		a.running = false
	}
}

func (a *agentT) configure(m *messaging.Message) {
	switch m.ContentType() {
	case private.ContentTypeConfiguration:
		cfg, status := private.ConfigurationContent(m)
		if !status.OK() {
			messaging.Reply(m, status, a.Name())
			return
		}
		a.state.Update(cfg)
	}
	messaging.Reply(m, messaging.StatusOK(), a.Name())
}

// Run - run the agent
func (a *agentT) run() {
}

func (a *agentT) exec(ctx context.Context, sql string, args ...any) (tag CommandTag, status *messaging.Status) {
	// Transaction processing.
	if a.state.DbClient == nil {
		return tag, messaging.NewStatus(messaging.StatusInvalidArgument, errors.New("DbClient is nil"))
	}
	txn, err0 := a.state.DbClient.Begin(ctx)
	if err0 != nil {
		return tag, messaging.NewStatus(StatusTxnBeginError, err0)
	}
	// Rollback is safe to call even if the tx is already closed, so if
	// the tx commits successfully, this is a no-op
	defer txn.Rollback(ctx)
	cmd, err := a.state.DbClient.Exec(ctx, sql, args)
	if err != nil {
		return newCmdTag(cmd), messaging.NewStatus(messaging.StatusInvalidArgument, recast(err))
	}
	err = txn.Commit(ctx)
	if err != nil {
		status = messaging.NewStatus(StatusTxnCommitError, err)
	} else {
		status = messaging.StatusOK()
	}
	return newCmdTag(cmd), status
}

func (a *agentT) statusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	return http.StatusInternalServerError
}

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

func (a *agentT) log(start time.Time, duration time.Duration, h http.Header, req *request, statusCode int) {
	if a.state.Log == nil {
		return
	}

	resp := newResponse(statusCode, nil)
	// TODO: determine how to set timeout from error
	if h != nil && h.Get(private.ThresholdRequest) != "" {
		// TODO : set the
		resp.Header().Set(private.ThresholdTimeoutName, "")
	}
	a.state.Log(private.TrafficEgress, start, duration, req.routeName, req, resp)
}
