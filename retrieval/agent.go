package retrieval

import (
	"context"
	"errors"
	"github.com/appellative-ai/core/messaging"
	"github.com/appellative-ai/postgres/private"
	"github.com/jackc/pgx/v5"
	"net/http"
	"time"
)

const (
	NamespaceName   = "sql:postgres:agent/retrieval"
	defaultDuration = time.Second * 3
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

func (a *agentT) retrieve(ctx context.Context, name, sql string, args []any) (rows pgx.Rows, err error) {
	if a.state.DbClient == nil {
		return nil, errors.New("DbClient is nil")
	}
	return a.state.DbClient.Query(ctx, sql, args)
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

func (a *agentT) log(start time.Time, duration time.Duration, req *request, statusCode int) {
	if a.state.Log == nil {
		return
	}

	resp := newResponse(statusCode)
	// TODO : set timeout value for the threshold header
	resp.Header().Set(private.ThresholdTimeoutName, "")
	a.state.Log(private.TrafficEgress, start, duration, req.routeName, req, resp)
}
