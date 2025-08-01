package retrieval

import (
	"bytes"
	"context"
	"time"
)

const (
	retrievalRouteName = "postgresql-retrieval"
)

type ScanFunc func(columnNames []string, values []any) error

type Interface struct {
	Marshal func(ctx context.Context, name, sql string, args ...any) (bytes.Buffer, error)
	Scan    func(ctx context.Context, fn ScanFunc, name, sql string, args ...any) error
}

// Retriever -
var Retriever = func() *Interface {
	return &Interface{
		Marshal: func(ctx context.Context, name, sql string, args ...any) (bytes.Buffer, error) {
			if ctx == nil {
				ctx = context.Background()
			}
			start := time.Now().UTC()
			rows, err := agent.retrieve(ctx, name, sql, args)
			agent.log(start, time.Since(start), retrievalRouteName, newRequest(name), newResponse(agent.statusCode(err)), ctx)
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
			agent.log(start, time.Since(start), retrievalRouteName, newRequest(name), newResponse(agent.statusCode(err)), ctx)
			if err != nil {
				return err
			}
			return Scanner(fn, createColumnNames(rows.FieldDescriptions()), rows)
		},
	}
}()
