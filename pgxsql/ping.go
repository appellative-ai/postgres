package pgxsql

import (
	"context"
	"errors"
)

// Ping - function for pinging the database cluster
func ping(ctx context.Context, req *request) (status error) {
	if dbClient == nil {
		return errors.New("error on PostgreSQL ping call : dbClient is nil")
	}
	ctx = req.setTimeout(ctx)
	err := dbClient.Ping(ctx)
	if err != nil {
		status = err
	} else {
		status = nil
	}
	return
}
