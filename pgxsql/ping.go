package pgxsql

import (
	"context"
	"errors"
	"github.com/behavioral-ai/core/messaging"
	"net/http"
)

// Ping - function for pinging the database cluster
func ping(ctx context.Context, req *request) (status *messaging.Status) {
	if dbClient == nil {
		return messaging.NewStatus(messaging.StatusInvalidArgument, errors.New("error on PostgreSQL ping call : dbClient is nil"))
	}
	ctx = req.setTimeout(ctx)
	err := dbClient.Ping(ctx)
	if err != nil {
		status = messaging.NewStatus(http.StatusInternalServerError, err)
	} else {
		status = messaging.StatusOK()
	}
	return
}
