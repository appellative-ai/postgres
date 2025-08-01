package retrieval

import (
	"fmt"
	"net/http"
)

const (
	protocol       = "message-based"
	postgresScheme = "postgres"
	queryMethod    = "query"
)

// Request - contains data needed to build the SQL statement related to the uri
type request struct {
	uri string
	h   http.Header
}

func newRequest(name string) *request {
	r := new(request)
	r.uri = name //buildUri(queryRoot, name)
	r.h = make(http.Header)
	return r
}

func (r *request) Header() http.Header {
	return r.h
}

func (r *request) Method() string {
	return queryMethod
}

func (r *request) Url() string {
	return r.uri
}

func (r *request) Protocol() string {
	return protocol
}

func buildUri(root, resource string) string {
	return fmt.Sprintf("%v://%v/%v:%v/%v/%v", postgresScheme, "host-name", "invalid-domain", "database-name", root, resource)
}
