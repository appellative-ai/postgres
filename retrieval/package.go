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
	Marshal func(ctx context.Context, name, sql string, args ...any) (bytes.Buffer, error)
	Scan    func(ctx context.Context, fn ScanFunc, name, sql string, args ...any) error
}

// Relation -
var Relation = func() *Resolution {
	return &Resolution{
		Marshal: func(ctx context.Context, name, sql string, args ...any) (bytes.Buffer, error) {
			newCtx, cancel := agent.setTimeout(ctx)
			defer cancel()

			start := time.Now().UTC()
			rows, err := agent.retrieve(newCtx, name, sql, args)
			agent.log(start, time.Since(start), newRequest(name, "template"), agent.statusCode(err))
			if err != nil {
				return bytes.Buffer{}, err
			}
			return Marshaler(createColumnNames(rows.FieldDescriptions()), rows)
		},
		Scan: func(ctx context.Context, fn ScanFunc, name, sql string, args ...any) error {
			newCtx, cancel := agent.setTimeout(ctx)
			defer cancel()

			start := time.Now().UTC()
			rows, err := agent.retrieve(newCtx, name, sql, args)
			agent.log(start, time.Since(start), newRequest(name, "template"), agent.statusCode(err))
			if err != nil {
				return err
			}
			return Scanner(fn, createColumnNames(rows.FieldDescriptions()), rows)
		},
	}
}()
