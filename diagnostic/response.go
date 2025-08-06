package diagnostic

import (
	"net/http"
)

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
