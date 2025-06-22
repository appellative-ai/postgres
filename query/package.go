package query

import (
	"context"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/postgres/common"
	"net/http"
	"time"
)

type RedirectFunc func() ([]byte, *messaging.Status)

// QueryT -  process a SQL select statement, returning a typed array
func QueryT[T common.Scanner[T]](ctx context.Context, h http.Header, resource, sql string, args ...any) (rows []T, status *messaging.Status) {
	newCtx, cancel := agent.setTimeout(ctx)
	defer cancel()

	path, code, ok := common.LocationValues(h)
	start := time.Now().UTC()
	if ok {
		status = messaging.NewStatus(code, nil)
		if path != "" {
			rows, status = common.Unmarshal[T](path)
		}
		agent.log(start, time.Since(start), h, newRequest(resource, "template"), status.Code)
		return
	}
	pgxRows, err := agent.query(newCtx, sql, args)
	statusCode := agent.statusCode(err)
	agent.log(start, time.Since(start), h, newRequest(resource, "template"), statusCode)
	if err != nil {
		return rows, messaging.NewStatus(statusCode, err)
	}
	return common.Scan[T](pgxRows)
}

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
