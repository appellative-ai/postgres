package request

import (
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
)

func newResult(tag pgconn.CommandTag) Result {
	return Result{
		Sql:          tag.String(),
		RowsAffected: tag.RowsAffected(),
		Insert:       tag.Insert(),
		Update:       tag.Update(),
		Delete:       tag.Delete(),
		Select:       tag.Select(),
	}
}

type response struct {
	statusCode int
	header     http.Header
}

func newResponse(statusCode int) *response {
	r := new(response)
	r.statusCode = statusCode
	r.header = make(http.Header)
	return r
}

func (r *response) StatusCode() int {
	return r.statusCode
}

func (r *response) Header() http.Header {
	return r.header
}
