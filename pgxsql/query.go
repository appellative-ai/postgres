package pgxsql

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
)

// Query - function for a Query
func query(ctx context.Context, req *request) (rows pgx.Rows, status error) {
	if req == nil {
		return nil, errors.New("error on PostgreSQL database retrieval call : request is nil")
	}
	if dbClient == nil {
		status = errors.New("error on PostgreSQL database retrieval call: dbClient is nil")
		return
	}
	var err error

	ctx = req.setTimeout(ctx)
	rows, err = dbClient.Query(ctx, buildSql(req), req.args)
	if err != nil {
		status = recast(err)
	} else {
		status = nil
	}
	return rows, status
}

// Scrap
//url, override := lookup.Value(req.test)
//defer apply(ctx, &newCtx, req, access.StatusCode(&status))
//if override {
//	// TO DO : create rows from file
//	return io2.New[pgx.Rows](url, nil)
//}
//var limited = false
//var fn func()
//
//fn, ctx, limited = controllerApply(ctx, startup.NewStatusCode(&status), req.Uri, messaging.ContextRequestId(ctx), "GET")
//defer fn()
//if limited {
//	return nil, messaging.NewStatus(messaging.StatusRateLimited)
//}
