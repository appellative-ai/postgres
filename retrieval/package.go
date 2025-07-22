package retrieval

import (
	"bytes"
	"context"
	"time"
)

// Sync with Access agent
const (
	egressTraffic = "egress"
	thresholdName = "x-threshold"
	timeoutName   = "timeout"
)

type ScanFunc func(columnNames []string, values []any) error

type Resolution struct {
	Marshal func(ctx context.Context, name, sql string, args ...any) (bytes.Buffer, error)
	Scan    func(ctx context.Context, fn ScanFunc, name, sql string, args ...any) error
}

// Relation -
var Relation = func() *Resolution {
	return &Resolution{
		Marshal: func(ctx context.Context, name, sql string, args ...any) (bytes.Buffer, error) {
			if ctx == nil {
				ctx = context.Background()
			}
			start := time.Now().UTC()
			rows, err := agent.retrieve(ctx, name, sql, args)
			agent.log(start, time.Since(start), newRequest(name, "template"), newResponse(agent.statusCode(err)), ctx)
			if err != nil {
				return bytes.Buffer{}, err
			}
			return Marshaler(createColumnNames(rows.FieldDescriptions()), rows)
		},
		Scan: func(ctx context.Context, fn ScanFunc, name, sql string, args ...any) error {
			if ctx == nil {
				ctx = context.Background()
			}
			start := time.Now().UTC()
			rows, err := agent.retrieve(ctx, name, sql, args)
			agent.log(start, time.Since(start), newRequest(name, "template"), newResponse(agent.statusCode(err)), ctx)
			if err != nil {
				return err
			}
			return Scanner(fn, createColumnNames(rows.FieldDescriptions()), rows)
		},
	}
}()
