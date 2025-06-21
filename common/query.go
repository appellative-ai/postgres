package common

import (
	"context"
	"fmt"
	"github.com/behavioral-ai/core/messaging"
	"net/http"
)

// QueryT -  process a SQL select statement, returning a typed array
func QueryT[T Scanner[T]](ctx context.Context, h http.Header, resource, sql string, args ...any) (rows []T, status *messaging.Status) {
	var t T

	if v, ok := any(t).(Variant); ok {
		if v != nil {
		}
		fmt.Printf("test: Query() -> [variant:%v]\n", ok)
	}
	rows, status = Scan[T](nil)
	return
}

/*
	//buf, status1, ok := processRedirect(args)
	start := time.Now().UTC()
	//if ok {
	//	status = status1
		if len(buf) >= 0 {
			rows, status = common.Unmarshal[T](buf)
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
*/
