package request

import (
	"context"
	"errors"
	"github.com/appellative-ai/postgres/common"
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

type Stat struct {
}

type Interface struct {
	Insert func(ctx context.Context, name, sql string, args ...any) (Response, error)
	Update func(ctx context.Context, name, sql string, args ...any) (Response, error)
	Delete func(ctx context.Context, name, sql string, args ...any) (Response, error)
	Relate func(ctx context.Context, instance, pattern, name, sql string, args ...any) (Response, error)
	Ping   func(ctx context.Context) error
	Stat   func() error
}

// Requester -
var Requester = func() *Interface {
	return &Interface{
		Insert: func(ctx context.Context, name, sql string, args ...any) (Response, error) {
			if ctx == nil {
				ctx = context.Background()
			}
			start := time.Now().UTC()

			/* TODO: determine how to bulk insert rows
			rows, status1 := common.Rows[T](entries)
			if !status1.OK() {
				agent.log(start, time.Since(start), h, newInsertRequest(name,sql,nil), status1.Code)
				return Response{}, status1
			}
			*/
			tag, err := agent.exec(ctx, sql, args)
			agent.log(start, time.Since(start), newInsertRequest(name, "", nil), newLogResponse(agent.statusCode(err)), ctx)
			return tag, err
		},
		Update: func(ctx context.Context, name, sql string, args ...any) (Response, error) {
			if ctx == nil {
				ctx = context.Background()
			}
			start := time.Now().UTC()
			tag, err := agent.exec(ctx, sql, args)
			agent.log(start, time.Since(start), newUpdateRequest(name, "", nil, nil), newLogResponse(agent.statusCode(err)), ctx)
			return tag, err

		},
		Delete: func(ctx context.Context, name, sql string, args ...any) (Response, error) {
			if ctx == nil {
				ctx = context.Background()
			}
			start := time.Now().UTC()
			tag, err := agent.exec(ctx, sql, args)
			agent.log(start, time.Since(start), newDeleteRequest(name, "", nil), newLogResponse(agent.statusCode(err)), ctx)
			return tag, err
		},
		Relate: func(ctx context.Context, instance, pattern, name, sql string, args ...any) (Response, error) {
			if ctx == nil {
				ctx = context.Background()
			}
			start := time.Now().UTC()
			resp, err := agent.relate(ctx, instance, pattern, name)
			agent.log(start, time.Since(start), newDeleteRequest("", "", nil), newLogResponse(agent.statusCode(err)), ctx)
			if err != nil {
				return resp, err
			}
			return resp, errors.New("not implemented")
		},
		Ping: func(ctx context.Context) error {
			if ctx == nil {
				ctx = context.Background()
			}
			start := time.Now().UTC()
			err := agent.ping(ctx)
			agent.log(start, time.Since(start), newDeleteRequest("", "", nil), newLogResponse(agent.statusCode(err)), ctx)
			return err
		},
		Stat: func() error {
			return agent.stat()
		},
	}
}()

// InsertT - execute a SQL insert statement
func InsertT[T common.Scanner[T]](ctx context.Context, h http.Header, resource, sql string, args ...any) (Response, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	count, status, ok := execValues(h)
	req := newInsertRequest(resource, "", nil, args...)

	start := time.Now().UTC()
	if ok {
		agent.log(start, time.Since(start), req, newLogResponse(status), ctx)
		return Response{RowsAffected: int64(count), Insert: true}, nil //messaging.NewStatus(status, nil)
	}
	/*
		 TODO: determine how to bulk insert rows
		rows, status1 := common.Rows[T](entries)
		if !status1.OK() {
			agent.log(start, time.Since(start), h, newInsertRequest(resource,sql,nil), status1.Code)
			return Response{}, status1
		}


	*/
	tag, err := agent.exec(ctx, sql, args)
	agent.log(start, time.Since(start), req, newLogResponse(agent.statusCode(err)), ctx)
	return tag, err
}

/*
// Update - execute a SQL update statement
func Update(ctx context.Context, h http.Header, resource, sql string, args ...any) (Response, *messaging.Status) {
	if ctx == nil {
		ctx = context.Background()
	}
	req := newUpdateRequest(resource, "", nil, nil)
	count, status, ok := execValues(h)
	start := time.Now().UTC()
	if ok {
		agent.log(start, time.Since(start), req, newLogResponse(status), ctx)
		return Response{RowsAffected: int64(count), Update: true}, messaging.NewStatus(status, nil)
	}
	tag, status1 := agent.exec(ctx, sql, args)
	agent.log(start, time.Since(start), req, newLogResponse(status1.Code), ctx)
	return tag, status1
}

// Delete - execute a SQL delete statement
func Delete(ctx context.Context, h http.Header, resource, sql string, args ...any) (Response, *messaging.Status) {
	if ctx == nil {
		ctx = context.Background()
	}
	count, status, ok := execValues(h)
	req := newDeleteRequest(resource, "", nil, args...)
	start := time.Now().UTC()
	if ok {
		agent.log(start, time.Since(start), req, newLogResponse(status), ctx)
		return Response{RowsAffected: int64(count), Delete: true}, messaging.NewStatus(status, nil)
	}
	tag, status1 := agent.exec(ctx, sql, args)
	agent.log(start, time.Since(start), req, newLogResponse(status1.Code), ctx)
	return tag, status1
}

*/
