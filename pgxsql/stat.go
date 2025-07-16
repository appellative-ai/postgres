package pgxsql

import (
	"errors"
	"github.com/appellative-ai/core/messaging"
	"github.com/jackc/pgx/v5/pgxpool"
)

func stat() (*pgxpool.Stat, *messaging.Status) {
	if dbClient == nil {
		return nil, messaging.NewStatus(messaging.StatusInvalidArgument, errors.New("error on PostgreSQL stat call : dbClient is nil"))
	}
	return dbClient.Stat(), messaging.StatusOK()
}

// Scrap
//var limited = false
//var fn func()

//fn, ctx, limited = controllerApply(ctx, startup.NewStatusCode(&status), StatUri, messaging.ContextRequestId(ctx), "GET")
//defer fn()
//if limited {
//	return nil, messaging.NewStatus(messaging.StatusRateLimited).SetRequestId(ctx)
//}
