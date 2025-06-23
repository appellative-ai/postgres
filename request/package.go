package request

import (
	"context"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/postgres/common"
	"net/http"
	"time"
)

// Response - results of a command
type Response struct {
	Sql          string `json:"sql"`
	RowsAffected int64  `json:"rows-affected"`
	Insert       bool   `json:"insert"`
	Update       bool   `json:"update"`
	Delete       bool   `json:"delete"`
	Select       bool   `json:"select"`
}

type Resolution struct {
	Insert func(ctx context.Context, resource, sql string, args ...any) (Response, *messaging.Status)
	Update func(ctx context.Context, resource, sql string, args ...any) (Response, *messaging.Status)
	Delete func(ctx context.Context, resource, sql string, args ...any) (Response, *messaging.Status)
}

// Mutation -
var Mutation = func() *Resolution {
	return &Resolution{
		Insert: func(ctx context.Context, resource, sql string, args ...any) (Response, *messaging.Status) {
			newCtx, cancel := agent.setTimeout(ctx)
			defer cancel()
			start := time.Now().UTC()

			/* TODO: determine how to bulk insert rows
			rows, status1 := common.Rows[T](entries)
			if !status1.OK() {
				agent.log(start, time.Since(start), h, newInsertRequest(resource,sql,nil), status1.Code)
				return Response{}, status1
			}
			*/
			tag, status1 := agent.exec(newCtx, sql, args)
			agent.log(start, time.Since(start), newInsertRequest(resource, "", nil), status1.Code)
			return tag, status1
		},
		Update: func(ctx context.Context, resource, sql string, args ...any) (Response, *messaging.Status) {
			newCtx, cancel := agent.setTimeout(ctx)
			defer cancel()
			start := time.Now().UTC()

			tag, status1 := agent.exec(newCtx, sql, args)
			agent.log(start, time.Since(start), newUpdateRequest(resource, "", nil, nil), status1.Code)
			return tag, status1

		},
		Delete: func(ctx context.Context, resource, sql string, args ...any) (Response, *messaging.Status) {
			newCtx, cancel := agent.setTimeout(ctx)
			defer cancel()

			start := time.Now().UTC()
			tag, status1 := agent.exec(newCtx, sql, args)
			agent.log(start, time.Since(start), newDeleteRequest(resource, "", nil), status1.Code)
			return tag, status1
		},
	}
}()

// InsertT - execute a SQL insert statement
func InsertT[T common.Scanner[T]](ctx context.Context, h http.Header, resource, sql string, args ...any) (Response, *messaging.Status) {
	newCtx, cancel := agent.setTimeout(ctx)
	defer cancel()

	count, status, ok := execValues(h)
	req := newInsertRequest(resource, "", nil, args...)

	start := time.Now().UTC()
	if ok {
		agent.log(start, time.Since(start), req, status)
		return Response{RowsAffected: int64(count), Insert: true}, messaging.NewStatus(status, nil)
	}
	/* TODO: determine how to bulk insert rows
	rows, status1 := common.Rows[T](entries)
	if !status1.OK() {
		agent.log(start, time.Since(start), h, newInsertRequest(resource,sql,nil), status1.Code)
		return Response{}, status1
	}
	*/
	tag, status1 := agent.exec(newCtx, sql, args)
	agent.log(start, time.Since(start), req, status1.Code)
	return tag, status1
}

// Update - execute a SQL update statement
func Update(ctx context.Context, h http.Header, resource, sql string, args ...any) (Response, *messaging.Status) {
	newCtx, cancel := agent.setTimeout(ctx)
	defer cancel()

	req := newUpdateRequest(resource, "", nil, nil)
	count, status, ok := execValues(h)
	start := time.Now().UTC()
	if ok {
		agent.log(start, time.Since(start), req, status)
		return Response{RowsAffected: int64(count), Update: true}, messaging.NewStatus(status, nil)
	}
	tag, status1 := agent.exec(newCtx, sql, args)
	agent.log(start, time.Since(start), req, status1.Code)
	return tag, status1
}

// Delete - execute a SQL delete statement
func Delete(ctx context.Context, h http.Header, resource, sql string, args ...any) (Response, *messaging.Status) {
	newCtx, cancel := agent.setTimeout(ctx)
	defer cancel()

	count, status, ok := execValues(h)
	req := newDeleteRequest(resource, "", nil, args...)
	start := time.Now().UTC()
	if ok {
		agent.log(start, time.Since(start), req, status)
		return Response{RowsAffected: int64(count), Delete: true}, messaging.NewStatus(status, nil)
	}
	tag, status1 := agent.exec(newCtx, sql, args)
	agent.log(start, time.Since(start), req, status1.Code)
	return tag, status1
}
