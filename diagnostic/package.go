package diagnostic

import (
	"context"
)

const (
	pingRouteName = "postgres-ping"
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
			return agent.ping(ctx)
		},
		Stat: func() error {
			return agent.stat()
		},
	}
}()
