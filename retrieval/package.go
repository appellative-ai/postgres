package retrieval

import (
	"bytes"
	"context"
	"time"
)

type Rows interface {
	Close()
	Next() bool
	Err() error
	Values() ([]any, error)
}

type ScanFunc func(columnNames []string, values []any) error

type Resolution struct {
	Marshal func(ctx context.Context, resource, sql string, args ...any) (bytes.Buffer, error)
	Scan    func(ctx context.Context, fn ScanFunc, resource, sql string, args ...any) error
}

// Relation -
var Relation = func() *Resolution {
	return &Resolution{
		Marshal: func(ctx context.Context, resource, sql string, args ...any) (bytes.Buffer, error) {
			newCtx, cancel := agent.setTimeout(ctx)
			defer cancel()

			start := time.Now().UTC()
			rows, err := agent.query(newCtx, sql, args)
			agent.log(start, time.Since(start), nil, newRequest(resource, "template"), agent.statusCode(err))
			if err != nil {
				return bytes.Buffer{}, err
			}
			return Marshaler(createColumnNames(rows.FieldDescriptions()), rows)
		},
		Scan: func(ctx context.Context, fn ScanFunc, resource, sql string, args ...any) error {
			newCtx, cancel := agent.setTimeout(ctx)
			defer cancel()

			start := time.Now().UTC()
			rows, err := agent.query(newCtx, sql, args)
			agent.log(start, time.Since(start), nil, newRequest(resource, "template"), agent.statusCode(err))
			if err != nil {
				return err
			}
			return Scanner(fn, createColumnNames(rows.FieldDescriptions()), rows)
		},
	}
}()
