package query

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

// Marshal -  process a SQL select statement, returning a JSON buffer
func Marshal(ctx context.Context, resource, sql string, args ...any) (bytes.Buffer, error) {
	newCtx, cancel := agent.setTimeout(ctx)
	defer cancel()

	start := time.Now().UTC()
	rows, err := agent.query(newCtx, sql, args)
	agent.log(start, time.Since(start), nil, newRequest(resource, "template"), agent.statusCode(err))
	if err != nil {
		return bytes.Buffer{}, err
	}
	return marshal(rows)
}

// Scan -  process a SQL select statement using a Scan function
func Scan(ctx context.Context, fn ScanFunc, resource, sql string, args ...any) error {
	newCtx, cancel := agent.setTimeout(ctx)
	defer cancel()

	start := time.Now().UTC()
	rows, err := agent.query(newCtx, sql, args)
	agent.log(start, time.Since(start), nil, newRequest(resource, "template"), agent.statusCode(err))
	if err != nil {
		return err
	}
	return Scanner(fn, createColumnNames(rows.FieldDescriptions()), rows)
}

/*
// QueryT -  process a SQL select statement, returning a typed array
func QueryT[T common.Scanner[T]](ctx context.Context, h http.Header, resource, sql string, args ...any) (rows []T, err error) {
	newCtx, cancel := agent.setTimeout(ctx)
	defer cancel()

	start := time.Now().UTC()
	pgxRows, err1 := agent.query(newCtx, sql, args)
	statusCode := agent.statusCode(err1)
	agent.log(start, time.Since(start), h, newRequest(resource, "template"), statusCode)
	if err1 != nil {
		return rows, err1
	}
	return common.Scan[T](pgxRows)
}


*/
/*
func processRedirect(args []any) ([]byte, *messaging.Status, bool) {
	if len(args) == 0 {
		return nil, nil, false
	}
	if fn, ok1 := args[0].(func() ([]byte, *messaging.Status)); ok1 {
		buf, status := fn()
		return buf, status, ok1
	}
	return nil, nil, false
}

*/

/*
	path, code, ok := queryValues(h)
	if ok {
		//status = messaging.NewStatus(code, nil)
		if path != "" {
			rows, status = common.Unmarshal[T](path)
		}
		agent.log(start, time.Since(start), h, newRequest(resource, "template"), status.Code)
		return
	}
*/

//type RedirectFunc func() ([]byte, *messaging.Status)
