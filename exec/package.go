package exec

import (
	"context"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/postgres/common"
	"net/http"
	"time"
)

// InsertT - execute a SQL insert statement
func InsertT[T common.Scanner[T]](ctx context.Context, h http.Header, resource, template string, entries []T, args ...any) (tag CommandTag, status *messaging.Status) {
	/* TODO : refactor
	_, _, stat1 := messaging.ExchangeHeaders(h)
	if stat1 != "" {
		start := time.Now().UTC()
		req := newInsertRequest(resource, template, nil, args...)
		ctx = req.setTimeout(ctx)
		status = jsonx.NewStatusFrom(stat1)
		log(start, h, req, status)
		return
	}
	*/
	rows, status1 := common.Rows[T](entries)
	if !status1.OK() {
		return CommandTag{}, status1
	}
	req := newInsertRequest(resource, template, rows, args...)
	start := time.Now().UTC()
	tag, status = exec(ctx, req)
	//log(start, h, req, status)
	return tag, status
}
