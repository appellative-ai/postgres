package request

import (
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
)

func newResponse(tag pgconn.CommandTag) Response {
	return Response{
		Sql:          tag.String(),
		RowsAffected: tag.RowsAffected(),
		Insert:       tag.Insert(),
		Update:       tag.Update(),
		Delete:       tag.Delete(),
		Select:       tag.Select(),
	}
}

type logResponse struct {
	statusCode int
	header     http.Header
}

func newLogResponse(statusCode int) *logResponse {
	r := new(logResponse)
	r.statusCode = statusCode
	r.header = make(http.Header)
	return r
}

func (r *logResponse) StatusCode() int {
	return r.statusCode
}

func (r *logResponse) Header() http.Header {
	return r.header
}
