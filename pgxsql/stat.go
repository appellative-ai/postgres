package pgxsql

import (
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func stat() (*pgxpool.Stat, error) {
	if dbClient == nil {
		return nil, errors.New("error on PostgreSQL stat call : dbClient is nil")
	}
	return dbClient.Stat(), nil // messaging.StatusOK()
}

// Scrap
//var limited = false
//var fn func()

//fn, ctx, limited = controllerApply(ctx, startup.NewStatusCode(&status), StatUri, messaging.ContextRequestId(ctx), "GET")
//defer fn()
//if limited {
//	return nil, messaging.NewStatus(messaging.StatusRateLimited).SetRequestId(ctx)
//}
