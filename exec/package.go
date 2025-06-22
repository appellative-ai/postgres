package exec

import (
	"context"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/postgres/common"
	"net/http"
	"time"
)

// CommandTag - results from an Exec command
type CommandTag struct {
	Sql          string `json:"sql"`
	RowsAffected int64  `json:"rows-affected"`
	Insert       bool   `json:"insert"`
	Update       bool   `json:"update"`
	Delete       bool   `json:"delete"`
	Select       bool   `json:"select"`
}

// InsertT - execute a SQL insert statement
func InsertT[T common.Scanner[T]](ctx context.Context, h http.Header, resource, sql string, entries []T, args ...any) (CommandTag, *messaging.Status) {
	newCtx, cancel := agent.setTimeout(ctx)
	defer cancel()

	count, status, ok := execValues(h)
	req := newInsertRequest(resource, "", nil, args...)

	start := time.Now().UTC()
	if ok {
		agent.log(start, time.Since(start), h, req, status)
		return CommandTag{RowsAffected: int64(count), Insert: true}, messaging.NewStatus(status, nil)
	}
	/* TODO: determine how to bulk insert rows
	rows, status1 := common.Rows[T](entries)
	if !status1.OK() {
		agent.log(start, time.Since(start), h, newInsertRequest(resource,sql,nil), status1.Code)
		return CommandTag{}, status1
	}
	*/
	tag, status1 := agent.exec(newCtx, sql, args)
	agent.log(start, time.Since(start), h, req, status1.Code)
	return tag, status1
}

// Update - execute a SQL update statement
func Update(ctx context.Context, h http.Header, resource, sql string, args ...any) (CommandTag, *messaging.Status) {
	newCtx, cancel := agent.setTimeout(ctx)
	defer cancel()

	req := newUpdateRequest(resource, "", nil, nil)
	count, status, ok := execValues(h)
	start := time.Now().UTC()
	if ok {
		agent.log(start, time.Since(start), h, req, status)
		return CommandTag{RowsAffected: int64(count), Update: true}, messaging.NewStatus(status, nil)
	}
	tag, status1 := agent.exec(newCtx, sql, args)
	agent.log(start, time.Since(start), h, req, status1.Code)
	return tag, status1
}

// Delete - execute a SQL delete statement
func Delete(ctx context.Context, h http.Header, resource, sql string, args ...any) (CommandTag, *messaging.Status) {
	newCtx, cancel := agent.setTimeout(ctx)
	defer cancel()

	count, status, ok := execValues(h)
	req := newDeleteRequest(resource, "", nil, args...)
	start := time.Now().UTC()
	if ok {
		agent.log(start, time.Since(start), h, req, status)
		return CommandTag{RowsAffected: int64(count), Delete: true}, messaging.NewStatus(status, nil)
	}
	tag, status1 := agent.exec(newCtx, sql, args)
	agent.log(start, time.Since(start), h, req, status1.Code)
	return tag, status1
}
