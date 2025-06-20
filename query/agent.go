package query

import (
	"context"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/postgres/private"
	"github.com/jackc/pgx/v5"
	"net/http"
	"time"
)

const (
	NamespaceName   = "sql:postgres:agent/query"
	defaultDuration = time.Second * 3
)

var (
	agent    *agentT
	cancelFn = func() {}
)

type agentT struct {
	running bool
	state   *private.Configuration
}

func NewAgent() messaging.Agent {
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

func (a *agentT) query(ctx context.Context, sql string, args []any) (rows pgx.Rows, err error) {
	//ctx = a.setTimeout(ctx)
	return a.state.DbClient.Query(ctx, sql, args)
}

/*
	func (a *agentT) setTimeout1(ctx context.Context) context.Context {
		if ctx == nil {
			return context.Background()
		}
		if d, ok := ctx.Deadline(); ok {
			a.state.Until = time.Until(d)
		}
		return ctx
	}
*/
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
