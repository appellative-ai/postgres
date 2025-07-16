package pgxsql

import (
	"context"
	"errors"
	"github.com/appellative-ai/core/messaging"
)

const (
	StatusTxnBeginError    = int(102) // Transaction processing begin error
	StatusTxnRollbackError = int(103) // Transaction processing rollback error
	StatusTxnCommitError   = int(104) // Transaction processing commit error
	StatusExecError        = int(105) // Execution error, as in a database call
	StatusNotStarted       = int(98)  // Not started

)

func exec(ctx context.Context, req *request) (tag CommandTag, status *messaging.Status) {
	if req == nil {
		return tag, messaging.NewStatus(messaging.StatusInvalidArgument, errors.New("error on PostgreSQL request call : request is nil"))
	}
	if dbClient == nil {
		status = messaging.NewStatus(messaging.StatusInvalidArgument, errors.New("error on PostgreSQL request call : dbClient is nil"))
		return
	}
	ctx = req.setTimeout(ctx)

	// Transaction processing.
	txn, err0 := dbClient.Begin(ctx)
	if err0 != nil {
		status = messaging.NewStatus(StatusTxnBeginError, err0)
		return tag, status
	}
	// Rollback is safe to call even if the tx is already closed, so if
	// the tx commits successfully, this is a no-op
	defer txn.Rollback(ctx)
	cmd, err := dbClient.Exec(ctx, buildSql(req), req.args)
	if err != nil {
		status = messaging.NewStatus(messaging.StatusInvalidArgument, recast(err))
		return newCmdTag(cmd), status
	}
	err = txn.Commit(ctx)
	if err != nil {
		status = messaging.NewStatus(StatusTxnCommitError, err)
	} else {
		status = messaging.StatusOK()
	}
	return newCmdTag(cmd), messaging.StatusOK()
}

// scrap
//defer apply(ctx, &newCtx, req, access.StatusCode(&status))
//if override {
//	return io2.New[CommandTag](url, nil)
//}
//
//url, override := lookup.Value(req.test)
//var newCtx context.Context
