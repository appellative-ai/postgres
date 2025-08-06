package diagnostic

import (
	"context"
	"time"
)

const (
	routeName = "postgres-ping"
	method    = "ping"
	name
)

type Stat struct {
}

// Interface -
type Interface struct {
	Ping func(ctx context.Context) error
	Stat func() error
}

// Diagnostic -
var Diagnostic = func() *Interface {
	return &Interface{
		Ping: func(ctx context.Context) error {
			if ctx == nil {
				ctx = context.Background()
			}
			start := time.Now().UTC()
			err := agent.ping(ctx)
			agent.log(start, time.Since(start), routeName, newRequest(name, method), newResponse(agent.statusCode(err)), ctx)
			return err
		},
		Stat: func() error {
			return agent.stat()
		},
	}
}()
