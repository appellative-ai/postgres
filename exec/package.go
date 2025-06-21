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
	tag, status, ok := processOverride(args)
	req := newInsertRequest(resource, "", nil, args...)

	start := time.Now().UTC()
	if ok {
		agent.log(start, time.Since(start), h, req, status.Code)
		return tag, status
	}
	/* TODO: determine how to bulk insert rows
	rows, status1 := common.Rows[T](entries)
	if !status1.OK() {
		agent.log(start, time.Since(start), h, newInsertRequest(resource,sql,nil), status1.Code)
		return CommandTag{}, status1
	}
	*/
	tag, status = agent.exec(newCtx, sql, args)
	agent.log(start, time.Since(start), h, req, status.Code)
	return tag, status
}

// Update - execute a SQL update statement
func Update(ctx context.Context, h http.Header, resource, sql string, args ...any) (CommandTag, *messaging.Status) {
	newCtx, cancel := agent.setTimeout(ctx)
	defer cancel()
	tag, status, ok := processOverride(args)
	req := newUpdateRequest(resource, "", nil, nil)
	start := time.Now().UTC()

	if ok {
		agent.log(start, time.Since(start), h, req, status.Code)
		return tag, status
	}
	tag, status = agent.exec(newCtx, sql, args)
	agent.log(start, time.Since(start), h, req, status.Code)
	return tag, status
}

// Delete - execute a SQL delete statement
func Delete(ctx context.Context, h http.Header, resource, sql string, args ...any) (CommandTag, *messaging.Status) {
	newCtx, cancel := agent.setTimeout(ctx)
	defer cancel()
	tag, status, ok := processOverride(args)
	req := newDeleteRequest(resource, "", nil, args...)
	start := time.Now().UTC()

	if ok {
		agent.log(start, time.Since(start), h, req, status.Code)
		return tag, status
	}
	tag, status = agent.exec(newCtx, sql, args)
	agent.log(start, time.Since(start), h, req, status.Code)
	return tag, status
}

func processOverride(args []any) (CommandTag, *messaging.Status, bool) {
	if len(args) == 0 {
		return CommandTag{}, nil, false
	}
	if fn, ok1 := args[0].(func() (CommandTag, *messaging.Status)); ok1 {
		tag, status := fn()
		return tag, status, ok1
	}
	return CommandTag{}, nil, false
}
