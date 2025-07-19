package retrieval

import (
	"context"
	"fmt"
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

func (r *response) SetTimeout(ctx context.Context) *response {
	if d, ok := ctx.Deadline(); ok {
		r.header.Add(thresholdName, fmt.Sprintf("%v=%v", timeoutName, d))
	}
	return r
}
