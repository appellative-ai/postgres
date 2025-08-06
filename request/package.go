package request

import (
	"context"
	"time"
)

const (
	routeName = "postgres-request"
)

// Result - results of a command
type Result struct {
	Sql          string `json:"sql"`
	RowsAffected int64  `json:"rows-affected"`
	Insert       bool   `json:"insert"`
	Update       bool   `json:"update"`
	Delete       bool   `json:"delete"`
	Select       bool   `json:"select"`
}

type Stat struct {
}

// Interface -
// TODO : determine if a bulk insert is needed
type Interface struct {
	Execute func(ctx context.Context, name, sql string, args ...any) (Result, error)
}

// Requester -
var Requester = func() *Interface {
	return &Interface{
		Execute: func(ctx context.Context, name, sql string, args ...any) (Result, error) {
			if ctx == nil {
				ctx = context.Background()
			}
			start := time.Now().UTC()
			tag, err := agent.exec(ctx, sql, args)
			agent.log(start, time.Since(start), routeName, newExecRequest(name), newResponse(agent.statusCode(err)), ctx)
			return tag, err
		},
	}
}()
