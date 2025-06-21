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
func InsertT[T common.Scanner[T]](ctx context.Context, h http.Header, resource, sql string, entries []T, args ...any) (tag CommandTag, status *messaging.Status) {
	rows, status1 := common.Rows[T](entries)
	if !status1.OK() {
		return CommandTag{}, status1
	}
	req := newInsertRequest(resource, "", rows, args...)
	start := time.Now().UTC()
	tag, status = agent.exec(ctx, sql, args)
	agent.log(start, time.Since(start), h, req, status.Code)
	return tag, status
}

// Update - execute a SQL update statement
func Update(ctx context.Context, h http.Header, resource, sql string, args ...any) (tag CommandTag, status *messaging.Status) {
	req := newUpdateRequest(resource, "", nil, nil)
	start := time.Now().UTC()
	tag, status = agent.exec(ctx, sql, args)
	agent.log(start, time.Since(start), h, req, status.Code)
	return tag, status
}

// Delete - execute a SQL delete statement
func Delete(ctx context.Context, h http.Header, resource, sql string, args ...any) (tag CommandTag, status *messaging.Status) {
	req := newDeleteRequest(resource, "", nil, args...)
	start := time.Now().UTC()
	tag, status = agent.exec(ctx, sql, args)
	agent.log(start, time.Since(start), h, req, status.Code)
	return tag, status
}
