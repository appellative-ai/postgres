package query

import (
	"context"
	"github.com/behavioral-ai/core/messaging"
	"net/http"
	"time"
)

const (
	PkgPath       = "github/behavioral-ai/postgres/pgxsql"
	userConfigKey = "user"
	pswdConfigKey = "pswd"
	uriConfigKey  = "uri"

	InsertRouteName = "postgresql-insert"
	UpdateRouteName = "postgresql-update"
	DeleteRouteName = "postgresql-delete"
	PingRouteName   = "postgresql-ping"
)

// QueryT -  process a SQL select statement, returning a type array
func QueryT[T Scanner[T]](ctx context.Context, h http.Header, sql string, args ...any) (rows []T, status *messaging.Status) {
	newCtx, cancel := agent.setTimeout(ctx)
	defer cancel()

	start := time.Now().UTC()
	r, err := agent.query(newCtx, sql, args)
	statusCode := agent.statusCode(err)
	agent.log(start, time.Since(start), h, newRequest("resource", "template"), statusCode)
	if err != nil {
		return rows, messaging.NewStatus(statusCode, err)
	}
	return Scan[T](r)
}
