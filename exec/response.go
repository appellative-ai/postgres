package exec

import "net/http"

type response struct {
	statusCode int
	header     http.Header
}

func newResponse(statusCode int, h http.Header) *response {
	r := new(response)
	r.statusCode = statusCode
	if h == nil {
		r.header = make(http.Header)
	} else {
		r.header = h
	}
	return r
}

func (r *response) StatusCode() int {
	return r.statusCode
}

func (r *response) Header() http.Header {
	return r.header
}
