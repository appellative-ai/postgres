package pgxsql

import (
	"context"
	"errors"
	"github.com/appellative-ai/core/messaging"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"time"
)

const (
	PkgPath       = "github/appellative-ai/postgres/pgxsql"
	userConfigKey = "user"
	pswdConfigKey = "pswd"
	uriConfigKey  = "uri"

	QueryRouteName  = "postgresql-retrieval"
	InsertRouteName = "postgresql-insert"
	UpdateRouteName = "postgresql-update"
	DeleteRouteName = "postgresql-delete"
	PingRouteName   = "postgresql-ping"
)

// Attr - key value pair
type Attr struct {
	Key string
	Val any
}

// Readiness - package readiness
func Readiness() error {
	if isReady() {
		return nil //messaging.StatusOK()
	}
	return errors.New("not started") //messaging.NewStatus(StatusNotStarted, "")
}

// QueryFunc - type declaration
type QueryFunc func(context.Context, http.Header, string, string, map[string][]string, ...any) (pgx.Rows, error)

// Query -  process a SQL select statement
func Query(ctx context.Context, h http.Header, resource, template string, values map[string][]string, args ...any) (rows pgx.Rows, status error) {
	req := newQueryRequestFromValues(resource, template, values, args...)
	start := time.Now().UTC()
	rows, status = query(ctx, req)
	log(start, h, req, status)
	return rows, status
}

// QueryFuncT - type declaration
type QueryFuncT[T Scanner[T]] func(context.Context, http.Header, string, string, map[string][]string, ...any) ([]T, error)

// QueryT -  process a SQL select statement, returning a type
func QueryT[T Scanner[T]](ctx context.Context, h http.Header, resource, template string, values map[string][]string, args ...any) (rows []T, status error) {
	req := newQueryRequestFromValues(resource, template, values, args...)
	req.Header().Set(messaging.XTo, "")
	start := time.Now().UTC()
	/* TODO - refactor
	_, resp, status1 := messaging.ExchangeHeaders(h)
	if resp != "" || status1 != "" {
		status = messaging.StatusOK()
		ctx = req.setTimeout(ctx)
		if resp != "" {
			rows, status = Unmarshal[T](resp)
		}
		//if status1 != "" {
		//	status = json.NewStatusFrom(status1)
		//}
		log(start, h, req, status)
		return
	}

	*/
	r, status2 := query(ctx, req)
	log(start, h, req, status2)
	if status2 != nil {
		return nil, status2
	}
	return Scan[T](r)
}

// InsertFunc - type
type InsertFunc func(context.Context, http.Header, string, string, [][]any, ...any) (CommandTag, error)

// Insert - execute a SQL insert statement
func Insert(ctx context.Context, h http.Header, resource, template string, values [][]any, args ...any) (tag CommandTag, status error) {
	req := newInsertRequest(resource, template, values, args...)
	start := time.Now().UTC()
	tag, status = exec(ctx, req)
	log(start, h, req, status)
	return tag, status
}

// InsertFuncT - type
type InsertFuncT[T Scanner[T]] func(context.Context, http.Header, string, string, []T, ...any) (CommandTag, error)

// InsertT - execute a SQL insert statement
func InsertT[T Scanner[T]](ctx context.Context, h http.Header, resource, template string, entries []T, args ...any) (tag CommandTag, status error) {
	/* TODO : refactor
	_, _, stat1 := messaging.ExchangeHeaders(h)
	if stat1 != "" {
		start := time.Now().UTC()
		req := newInsertRequest(resource, template, nil, args...)
		ctx = req.setTimeout(ctx)
		status = jsonx.NewStatusFrom(stat1)
		log(start, h, req, status)
		return
	}
	*/
	rows, status1 := Rows[T](entries)
	if status1 != nil {
		return CommandTag{}, status1
	}
	req := newInsertRequest(resource, template, rows, args...)
	start := time.Now().UTC()
	tag, status = exec(ctx, req)
	log(start, h, req, status)
	return tag, status
}

// UpdateFunc - type
type UpdateFunc func(context.Context, http.Header, string, string, []Attr, []Attr) (CommandTag, error)

// Update - execute a SQL update statement
func Update(ctx context.Context, h http.Header, resource, template string, where []Attr, args []Attr) (tag CommandTag, status error) {
	req := newUpdateRequest(resource, template, convert(where), convert(args))
	start := time.Now().UTC()
	tag, status = exec(ctx, req)
	log(start, h, req, status)
	return tag, status
}

// DeleteFunc - type
type DeleteFunc func(context.Context, http.Header, string, string, []Attr, ...any) (CommandTag, error)

// Delete - execute a SQL delete statement
func Delete(ctx context.Context, h http.Header, resource, template string, where []Attr, args ...any) (tag CommandTag, status error) {
	req := newDeleteRequest(resource, template, convert(where), args...)
	start := time.Now().UTC()
	tag, status = exec(ctx, req)
	log(start, h, req, status)
	return tag, status
}

// Stat - retrieve Pgx pool stats
func Stat() (*pgxpool.Stat, error) {
	return stat()
}

// Ping - ping the database cluster
func Ping(ctx context.Context, h http.Header) error {
	req := newPingRequest()
	start := time.Now().UTC()
	status := ping(ctx, req)
	log(start, h, req, status)
	return status
}
